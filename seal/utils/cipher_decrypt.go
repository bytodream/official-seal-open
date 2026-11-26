// cipher_decrypt.go 提供解密/还原方向的密钥解析与混合解密还原能力。
//
// 主要能力:
//   - 密钥 PEM 文件读取与解析（支持口令加密私钥、PKIX/PKCS#1 公钥、X.509 证书）
//   - HybridDecrypt: 用 RSA 私钥解开 HybridEncrypt 载荷（OAEP-SHA256 + AES-256-GCM）
//   - HybridOpen: 用 RSA 公钥还原 HybridSeal 载荷（PKCS#1 v1.5 + AES-256-GCM）
//
// 加密/封装方向的实现见 cipher_encrypt.go。
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strings"
)

// 混合加密二进制载荷格式（版本 1 与 2 布局相同，仅 RSA 封装数据密钥的方式不同）：
//
//	[1 字节版本][2 字节大端 keyBlob 长度][RSA 封装的数据密钥][12 字节 GCM nonce][AES-GCM 密文含 tag]
//
// 版本 1：数据密钥由公钥 RSA-OAEP(SHA-256) 封装，私钥解开。
// 版本 2：数据密钥由私钥 PKCS#1 v1.5 封装，公钥解开。
// 整段载荷再 Base64 编码为文本，便于环境变量等渠道传递。
const (
	defaultRSAKeyBits          = 2048 // defaultRSAKeyBits 默认 RSA 密钥长度。
	hybridCipherVersionEncrypt = 1    // 公钥 OAEP 封装 AES 密钥
	hybridCipherVersionSeal    = 2    // 私钥 PKCS#1 v1.5 封装 AES 密钥（公钥解开）
	hybridAESKeySize           = 32   // AES-256 数据密钥长度（字节）
	hybridGCMNonceSize         = 12   // GCM nonce 长度（字节）
)

// PEM 块类型。
const (
	pemTypeRSAPrivateKey = "RSA PRIVATE KEY" // PKCS#1 私钥
	pemTypeRSAPublicKey  = "RSA PUBLIC KEY"  // PKCS#1 公钥
	pemTypePublicKey     = "PUBLIC KEY"      // PKIX 公钥（SubjectPublicKeyInfo）
	pemTypeCertificate   = "CERTIFICATE"     // X.509 证书（可从中提取公钥）
)

// ---- 密钥文件读取与解析 ----

// GetPrivateKeyByPath 读取 path 处的 PEM 私钥并解析。
// 私钥被口令加密时必须传入正确的 password。
func GetPrivateKeyByPath(path string, password string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key %s: %w", path, err)
	}
	return GetPrivateKeyByBytes(keyBytes, password)
}

// GetPrivateKeyByBytes 解析 PEM 编码的 PKCS#1 私钥。
// 支持口令加密的私钥（传统加密 PEM 格式）；password 为空时按明文解析。
func GetPrivateKeyByBytes(privateKeyBytes []byte, password string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || len(block.Bytes) == 0 {
		return nil, fmt.Errorf("invalid private key PEM")
	}

	derBytes := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		if password == "" {
			return nil, fmt.Errorf("private key is password protected, password required")
		}
		decrypted, err := x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("decrypt private key PEM: %w", err)
		}
		derBytes = decrypted
	}
	key, err := x509.ParsePKCS1PrivateKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS#1 private key: %w", err)
	}
	return key, nil
}

// GetPublicKeyByPath 读取 path 处的 PEM 文件并提取 RSA 公钥。
// 支持 PKIX 公钥（PUBLIC KEY）、PKCS#1 公钥（RSA PUBLIC KEY）与 X.509 证书（CERTIFICATE）。
func GetPublicKeyByPath(path string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read public key %s: %w", path, err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || len(block.Bytes) == 0 {
		return nil, fmt.Errorf("invalid public key PEM in %s", path)
	}
	switch block.Type {
	case pemTypePublicKey, pemTypeRSAPublicKey:
		return parseRSAPublicKeyFromPEMBlock(block)
	case pemTypeCertificate:
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse certificate: %w", err)
		}
		rsaPub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate public key is not RSA: %T", cert.PublicKey)
		}
		return rsaPub, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}

// GetPublicKeyByPrivateKeyPath 读取私钥文件并推导出对应的公钥
// （RSA 私钥中包含公钥的全部信息，可无损还原）。
func GetPublicKeyByPrivateKeyPath(privateKeyPath string, password string) (*rsa.PublicKey, error) {
	privateKey, err := GetPrivateKeyByPath(privateKeyPath, password)
	if err != nil {
		return nil, err
	}
	return &privateKey.PublicKey, nil
}

// GetPublicKeyByPrivateKeyBytes 解析私钥 PEM 字节并直接返回对应的 RSA 公钥。
func GetPublicKeyByPrivateKeyBytes(privateKeyBytes []byte, password string) (*rsa.PublicKey, error) {
	privateKey, err := GetPrivateKeyByBytes(privateKeyBytes, password)
	if err != nil {
		return nil, err
	}
	return &privateKey.PublicKey, nil
}

// parsePublicKeyPEM 从 PEM 字符串解析 RSA 公钥。
func parsePublicKeyPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || len(block.Bytes) == 0 {
		return nil, fmt.Errorf("invalid public key PEM")
	}
	return parseRSAPublicKeyFromPEMBlock(block)
}

// parseRSAPublicKeyFromPEMBlock 解析 PEM 块中的 RSA 公钥，
// 支持 PKIX（PUBLIC KEY）与 PKCS#1（RSA PUBLIC KEY）两种块类型。
func parseRSAPublicKeyFromPEMBlock(block *pem.Block) (*rsa.PublicKey, error) {
	switch block.Type {
	case pemTypePublicKey:
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKIX public key: %w", err)
		}
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key: %T", pub)
		}
		return rsaPub, nil
	case pemTypeRSAPublicKey:
		rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKCS#1 public key: %w", err)
		}
		return rsaPub, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}

// ---- 混合加密解密：私钥解开 HybridEncrypt 载荷（AES-256-GCM + RSA-OAEP-SHA256）----
// 加密方向（HybridEncrypt）见 cipher_encrypt.go。

// HybridDecrypt 使用 RSA 私钥解密 HybridEncrypt 返回的 Base64 载荷，还原明文。
// privateKeyPEM 为 PEM 字符串；私钥块带口令保护时需传入正确 password，明文私钥时 password 可为空。
func HybridDecrypt(privateKeyPEM, password, encoded string) ([]byte, error) {
	priv, err := GetPrivateKeyByBytes([]byte(privateKeyPEM), password)
	if err != nil {
		return nil, err
	}
	blob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid hybrid ciphertext base64: %w", err)
	}
	return hybridDecrypt(priv, blob)
}

// hybridDecrypt 解析版本 1 载荷并用私钥 OAEP 解开数据密钥，再 AES-GCM 解密正文。
func hybridDecrypt(priv *rsa.PrivateKey, blob []byte) ([]byte, error) {
	const headerSize = 1 + 2
	if len(blob) < headerSize+hybridGCMNonceSize {
		return nil, fmt.Errorf("invalid hybrid ciphertext: too short")
	}
	if blob[0] != hybridCipherVersionEncrypt {
		return nil, fmt.Errorf("unsupported hybrid cipher version: %d", blob[0])
	}
	encKeyLen := int(binary.BigEndian.Uint16(blob[1:3]))
	rest := blob[headerSize:]
	minBody := encKeyLen + hybridGCMNonceSize + aes.BlockSize
	if len(rest) < minBody {
		return nil, fmt.Errorf("invalid hybrid ciphertext: truncated body")
	}
	encKey := rest[:encKeyLen]
	nonce := rest[encKeyLen : encKeyLen+hybridGCMNonceSize]
	sealed := rest[encKeyLen+hybridGCMNonceSize:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encKey, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA-OAEP decrypt AES key: %w", err)
	}
	if len(aesKey) != hybridAESKeySize {
		return nil, fmt.Errorf("invalid decrypted AES key length")
	}
	return aesGCMOpen(aesKey, nonce, sealed)
}

// ---- 混合封装还原：公钥解开 HybridSeal 载荷（AES-256-GCM + RSA PKCS#1 v1.5）----
// 封装方向（HybridSeal）见 cipher_encrypt.go。

// HybridOpen 使用 RSA 公钥解密 HybridSeal 返回的 Base64 载荷（版本 2），还原明文。
func HybridOpen(publicKeyPEM, encoded string) ([]byte, error) {
	encoded = strings.TrimSpace(encoded)
	pub, err := parsePublicKeyPEM(publicKeyPEM)
	if err != nil {
		return nil, err
	}
	blob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid hybrid seal base64: %w", err)
	}
	return hybridOpen(pub, blob)
}

// hybridOpen 解析版本 2 载荷并用公钥解开数据密钥，再 AES-GCM 解密正文。
func hybridOpen(pub *rsa.PublicKey, blob []byte) ([]byte, error) {
	wrappedKey, nonce, sealed, err := parseHybridBlob(blob, hybridCipherVersionSeal)
	if err != nil {
		return nil, err
	}
	aesKey, err := rsaPublicDecryptPKCS1v15(pub, wrappedKey)
	if err != nil {
		return nil, fmt.Errorf("RSA public unwrap AES key: %w", err)
	}
	if len(aesKey) != hybridAESKeySize {
		return nil, fmt.Errorf("invalid unwrapped AES key length")
	}
	return aesGCMOpen(aesKey, nonce, sealed)
}

// aesGCMOpen 使用 AES-256-GCM 校验并解密 sealed，nonce 须与加密时一致。
func aesGCMOpen(aesKey, nonce, sealed []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("AES-GCM: %w", err)
	}
	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, fmt.Errorf("AES-GCM decrypt: %w", err)
	}
	return plaintext, nil
}

// parseHybridBlob 从二进制载荷中解析出 RSA 封装的密钥、GCM nonce 与密文；wantVersion 须与首字节版本一致。
func parseHybridBlob(blob []byte, wantVersion byte) (wrappedKey, nonce, sealed []byte, err error) {
	const headerSize = 1 + 2
	if len(blob) < headerSize+hybridGCMNonceSize {
		return nil, nil, nil, fmt.Errorf("invalid hybrid blob: too short")
	}
	if blob[0] != wantVersion {
		return nil, nil, nil, fmt.Errorf("unsupported hybrid cipher version: %d", blob[0])
	}
	wrappedLen := int(binary.BigEndian.Uint16(blob[1:3]))
	rest := blob[headerSize:]
	minBody := wrappedLen + hybridGCMNonceSize + aes.BlockSize
	if len(rest) < minBody {
		return nil, nil, nil, fmt.Errorf("invalid hybrid blob: truncated body")
	}
	wrappedKey = rest[:wrappedLen]
	nonce = rest[wrappedLen : wrappedLen+hybridGCMNonceSize]
	sealed = rest[wrappedLen+hybridGCMNonceSize:]
	return wrappedKey, nonce, sealed, nil
}

// rsaPublicDecryptPKCS1v15 用公钥还原经私钥 PKCS#1 v1.5 类型 1 填充的短消息（标准库未导出该运算）。
func rsaPublicDecryptPKCS1v15(pub *rsa.PublicKey, ciphertext []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if len(ciphertext) != k {
		return nil, fmt.Errorf("rsa: ciphertext length mismatch")
	}
	c := new(big.Int).SetBytes(ciphertext)
	e := big.NewInt(int64(pub.E))
	m := new(big.Int).Exp(c, e, pub.N)
	em := rsaPadToKeySize(m.Bytes(), k)
	if len(em) < 11 || em[0] != 0x00 || em[1] != 0x01 {
		return nil, fmt.Errorf("rsa: invalid padding")
	}
	i := 2
	for ; i < len(em); i++ {
		if em[i] == 0x00 {
			break
		}
		if em[i] != 0xff {
			return nil, fmt.Errorf("rsa: invalid padding")
		}
	}
	i++
	if i >= len(em) {
		return nil, fmt.Errorf("rsa: invalid padding")
	}
	return em[i:], nil
}

// rsaPadToKeySize 将大整数解码后的字节左补零至 RSA 模长（字节），供 PKCS#1 填充解析使用。
func rsaPadToKeySize(b []byte, size int) []byte {
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}
