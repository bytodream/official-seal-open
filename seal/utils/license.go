package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bytodream/official-seal-open/seal/model"
)

// 用于验证许可证签名的 RSA 公钥（PEM 格式）。
const licensePublicKeyPEM = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3kyTs/L6Ia8oe8LrKYHA\nDMIgsFerYRoIJbsptmev83ciGeMzo9HwSE9dAJS6ByCMmXYTk6idONMHffvuYxXb\nV9yi0wgYmBSgjMFTTbaEZrATJRIjvvrfmoi1CLNc3RO813NqnLK9TUAHLZ7FLFzp\nNQ8QLbmn5uDFfdI0xwVpmxCUdV+N9l0UMDrmxpbj5BD7LxP0pE1OmAaSYnnZaXob\nMVdJ99K8nSBuC5LrytCwFCZ3HFEmpEFWAn3GWZuQ5Jzoao9iS62EhTmFcG20DITg\njKD0YZi6mw1En1+oNwHJcz5kUfQv1cqhEV1D2qv14S3f1WXPQDdXHKgMnZwZhU6Z\nkQIDAQAB\n-----END PUBLIC KEY-----\n"

const communityLicense = "AgEAy7qZ7+eVbLwP0H7HR0MUd9IppoSqwvSgdr6IWTD8yzN5gesf5PgNMi4zkX58KTJB9Kk5F0YFtv7Pfr1TmoaDwv3ksIDMU63RNYEt8HAuM4j4r2b0kC52E96SgSmamG1aw3Gw2lDKDT+wdkzwl71ZuOx2b7prCcecTELdv9nJ+tyDJqYj79p1R1O0CzAHL1Ra1jmhzj7RVKwJNg5Ua7+cIzeyFp+pGfkTAuO0cHxoiT+ag5O+bnIgZdzSNMCAlLdrJt0s5qipt99dGzAjRTs4cDOBozMRE5YcFfG09IRjaMM4dpmSVTECHUonZ2RvyYpjNiLM2IsZo0Hzt1V26qUM6+97SwrB2tRnN9PveYzWZo5MKQRfyCWgYVQ+ORx8i48rg4iBbd0s0vjCwstLR83zeaE4/hnZL4uZzkdK3THT2lUSghfEYCdeGvXDkQF7rTsiNmUoJGEkU+j8wdo1c2f1I0zZ8fFFSYD0t6RbHlgagyRoHQF6145sOizJSgTXG68RP3EbabACcJ4FXglD7S/YFniDCg9tbwEchLUZyALALdERVyY7y8gbVYlfowJdeTvBA0Zvjp6SvrdjQK+5dPcncFTRN6vboghAcExbKDKYEWzn30Kt6lE7xMvxmjnnES5/LW/Py30mtGD8ULgFSmF0MWuTTTvTWuhW3awHfpA5Ku4pxmq2llTxjpd8EaAZyLfCTGbg"

// LicenseTier 许可证类型。
const (
	LicenseTierUnlicensed = "unlicensed" // 未授权
	LicenseTierCommunity  = "community"  // 开源社区版（AGPLv3）
	LicenseTierIndividual = "individual" // 商业版-个人
	LicenseTierBusiness   = "business"   // 商业版-企业
)

// 许可证与章面权限相关错误。
var (
	ErrLicenseRequired         = errors.New("valid license required for seal generation")
	ErrLicenseExpired          = errors.New("license expired")
	ErrSealCategoryNotLicensed = errors.New("seal category not allowed by license")
	ErrUsageLimitExceeded      = errors.New("license usage limit exceeded")
)

// init 从 OS_LICENSE（或兼容 OFFICIAL_SEAL_LICENSE）自动加载许可证。
func init() {
	license := os.Getenv("OS_LICENSE")
	if license == "" {
		license = os.Getenv("OFFICIAL_SEAL_LICENSE")
	}
	if license == "" {
		license = communityLicense
	}
	if err := SetLicenseKey(license); err != nil {
		Log.Error("Unable to load license code: %v", err)
	}
}

// License 描述一份许可证密钥。
type License struct {
	LicenseId    string     `json:"license_id"`    // 许可证唯一标识
	CustomerId   string     `json:"customer_id"`   // 客户唯一标识
	CustomerName string     `json:"customer_name"` // 客户名称（需与设置时传入的名称一致）
	Tier         string     `json:"tier"`          // 许可证类型，见 LicenseTier 常量
	CreatedAt    time.Time  `json:"-"`             // 创建时间（由 CreatedAtInt 解析而来）
	CreatedAtInt int64      `json:"created_at"`    // 创建时间的 Unix 时间戳
	ExpiresAt    *time.Time `json:"-"`             // 过期时间（由 ExpiresAtInt 解析而来），nil 表示永不过期
	ExpiresAtInt int64      `json:"expires_at"`    // 过期时间的 Unix 时间戳，0 表示永不过期
	CreatedBy    string     `json:"created_by"`    // 创建者标识
	CreatorName  string     `json:"creator_name"`  // 签发人姓名
	CreatorEmail string     `json:"creator_email"` // 签发人邮箱
	// SealCategories 允许的章面分类（model.SealCategory 字符串），是章型权限的唯一依据：
	// 由签发方显式配置，未配置时不允许任何章型，不随 Tier 推导默认值。
	SealCategories []string `json:"seal_categories,omitempty"`
	// UsageLimit 使用次数上限；0 表示不限制，大于 0 时每次生成章面都会计数，用尽后拒绝生成。
	// 已使用次数持久化在本地 usage.json 中，按 LicenseId 区分。
	UsageLimit int64 `json:"usage_limit,omitempty"`
}

// 全局许可证状态。
var (
	// currentLicenseKey 当前生效的许可证密钥，默认为"未授权"占位密钥。
	// 生成流程会并发读取，运行期 SetLicenseKey 可能并发写入，故用读写锁保护。
	currentLicenseKey = MakeUnlicensedKey()
	licenseMu         sync.RWMutex
)

// MakeUnlicensedKey 生成"未授权"信息。
func MakeUnlicensedKey() *License {
	key := License{
		CustomerName: "Unlicensed",
		Tier:         LicenseTierUnlicensed,
	}
	return &key
}

// Validate License.Validate 校验许可证有效性。
func (key *License) Validate() error {
	if !key.IsLicensed() {
		return fmt.Errorf("invalid license: Tier is invalid")
	}
	if ReleasedAt.After(key.CreatedAt) {
		return fmt.Errorf("invalid license: Created At is invalid")
	}
	if key.ExpiresAt == nil {
		expiresAt := key.CreatedAt.AddDate(1, 0, 0)
		if ReleasedAt.After(expiresAt) {
			expiresAt = ReleasedAt
		}
		key.ExpiresAt = &expiresAt
	}
	if key.ExpiresAtInt < 0 || key.UsageLimit < 0 {
		return fmt.Errorf("invalid license: Expires At / Usage Limit cannot be negative")
	}
	if key.CreatedAt.After(*key.ExpiresAt) {
		return fmt.Errorf("invalid license: Created At cannot be Greater than Expires At")
	}
	if key.isExpired() {
		message := "The license has already expired.\n" +
			"You may need to update the license key to the newest license key for your organization."
		return fmt.Errorf("invalid license: %s", message)
	}
	return nil
}

// isExpired 判断许可证是否已过期。
func (lk *License) isExpired() bool {
	if lk.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*lk.ExpiresAt)
}

// IsLicensed 判断是否已持有有效许可证。
func (lk *License) IsLicensed() bool {
	return lk.Tier != LicenseTierUnlicensed
}

// TypeToString 返回许可证类型的可读描述。
func (lk *License) TypeToString() string {
	switch {
	case lk.Tier == LicenseTierUnlicensed:
		return "Unlicensed"
	case lk.Tier == LicenseTierCommunity:
		return "AGPLv3 Open Source Community License"
	case lk.Tier == LicenseTierIndividual || lk.Tier == "indie":
		return "Commercial License - Individual"
	}
	return "Commercial License - Business"
}

// ToString 返回许可证的完整可读描述。
func (lk *License) ToString() string {
	desc := fmt.Sprintf("License Id: %s\n", lk.LicenseId)
	desc += fmt.Sprintf("Customer Id: %s\n", lk.CustomerId)
	desc += fmt.Sprintf("Customer Name: %s\n", lk.CustomerName)
	desc += fmt.Sprintf("Tier: %s\n", lk.Tier)
	desc += fmt.Sprintf("Created At: %s\n", UtcTimeFormat(lk.CreatedAt))
	if lk.ExpiresAt == nil {
		desc += "Expires At: Never\n"
	} else {
		desc += fmt.Sprintf("Expires At: %s\n", UtcTimeFormat(*lk.ExpiresAt))
	}
	if len(lk.SealCategories) > 0 {
		desc += fmt.Sprintf("Seal Categories: %s\n", strings.Join(lk.SealCategories, ", "))
	} else {
		desc += "Seal Categories: None\n"
	}
	if lk.UsageLimit > 0 {
		desc += fmt.Sprintf("Usage Limit: %d times\n", lk.UsageLimit)
	} else {
		desc += "Usage Limit: Unlimited\n"
	}
	desc += fmt.Sprintf("Creator: %s <%s>\n", lk.CreatorName, lk.CreatorEmail)
	return desc
}

// GetLicenseKey 返回当前生效许可证的副本。
func GetLicenseKey() *License {
	licenseMu.RLock()
	defer licenseMu.RUnlock()
	if currentLicenseKey == nil {
		return nil
	}
	key := *currentLicenseKey
	// 切片字段深拷贝，避免调用方改写与全局共享底层数组。
	key.SealCategories = append([]string(nil), currentLicenseKey.SealCategories...)
	return &key
}

// SetLicenseKey 解析并在本地校验许可证密钥（HybridOpen 还原正文后校验，不请求后端）。
func SetLicenseKey(content string) error {
	key, err := parseLicenseKey(content)
	if err != nil {
		Log.Error("License code decode error: %v", err)
		return err
	}
	if err = key.Validate(); err != nil {
		Log.Error("License code validation error: %v", err)
		return err
	}
	licenseMu.Lock()
	currentLicenseKey = &key
	licenseMu.Unlock()
	return nil
}

// EnsureLicenseActive 校验当前已加载的许可证存在且未过期（不包含章型判断）。
func EnsureLicenseActive() error {
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		return ErrLicenseRequired
	}
	if lk.isExpired() {
		return ErrLicenseExpired
	}
	return nil
}

// parseLicenseKey 从许可证密文（HybridSeal 载荷）中还原 JSON 并反序列化出 License。
func parseLicenseKey(content string) (License, error) {
	var key License
	decoded, err := HybridOpen(licensePublicKeyPEM, content)
	if err != nil {
		return key, err
	}
	if err := json.Unmarshal(decoded, &key); err != nil {
		return key, err
	}
	key.CreatedAt = time.Unix(key.CreatedAtInt, 0)
	if key.ExpiresAtInt > 0 {
		expiresAt := time.Unix(key.ExpiresAtInt, 0)
		key.ExpiresAt = &expiresAt
	}
	return key, nil
}

// ---- 章面类型权限 ----

// CheckCustomSealAttribute 根据自定义章面参数中的 SealCategory 校验许可证（nil 时视为 enterprise）。
func CheckCustomSealAttribute(a *model.CustomSealAttribute) error {
	if a == nil {
		return model.ErrSealAttributeNotConfigured
	}
	return CheckSealCategory(sealCategoryFromAttribute(a))
}

// sealCategoryFromAttribute 解析请求章型，未指定时默认为 enterprise。
func sealCategoryFromAttribute(a *model.CustomSealAttribute) string {
	if a.SealCategory == nil {
		return string(model.SealCategoryEnterprise)
	}
	return string(*a.SealCategory)
}

// CheckCustomSealAttributePermission 按自定义章面参数中的 SealCategory 校验许可证有效性与章型权限，
// 不检查使用次数（供预览等不消耗额度的路径使用）。
func CheckCustomSealAttributePermission(a *model.CustomSealAttribute) error {
	if a == nil {
		return model.ErrSealAttributeNotConfigured
	}
	return CheckSealCategoryPermission(sealCategoryFromAttribute(a))
}

// CheckSealCategoryPermission 只校验许可证有效性与章型权限，不检查使用次数。
func CheckSealCategoryPermission(category string) error {
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		return ErrLicenseRequired
	}
	if lk.isExpired() {
		return ErrLicenseExpired
	}
	cat := string(model.ConvertSealCategory(category))
	if !lk.allowsSealCategory(cat) {
		return fmt.Errorf("%w: %s", ErrSealCategoryNotLicensed, cat)
	}
	return nil
}

// CheckSealCategory 校验当前许可证是否允许生成指定章面分类（含使用次数检查）；category 与 model.SealCategory 取值一致。
func CheckSealCategory(category string) error {
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		return ErrLicenseRequired
	}
	if lk.isExpired() {
		return ErrLicenseExpired
	}
	if err := lk.checkUsageQuota(); err != nil {
		return err
	}
	cat := string(model.ConvertSealCategory(category))
	if !lk.allowsSealCategory(cat) {
		return fmt.Errorf("%w: %s", ErrSealCategoryNotLicensed, cat)
	}
	return nil
}

// AllowsSealCategory 判断许可证是否允许生成 category 章型。
func (lk *License) AllowsSealCategory(category string) bool {
	if lk == nil || !lk.IsLicensed() {
		return false
	}
	if lk.isExpired() {
		return false
	}
	if lk.usageExceeded() {
		return false
	}
	cat := string(model.ConvertSealCategory(category))
	return lk.allowsSealCategory(cat)
}

// AllowedSealCategories 返回当前许可证允许的章面分类列表（已规范化，按字典序）。
func (lk *License) AllowedSealCategories() []string {
	if lk == nil {
		return nil
	}
	set := lk.effectiveSealCategorySet()
	out := make([]string, 0, len(set))
	for c := range set {
		out = append(out, c)
	}
	sort.Strings(out)
	return out
}

func (lk *License) allowsSealCategory(category string) bool {
	_, ok := lk.effectiveSealCategorySet()[category]
	return ok
}

// effectiveSealCategorySet 仅使用许可证显式配置的 SealCategories（已规范化）。
// 章型权限完全由签发内容决定：未配置时不允许任何章型，不按 tier 推导默认值。
func (lk *License) effectiveSealCategorySet() map[string]struct{} {
	set := make(map[string]struct{}, len(lk.SealCategories))
	for _, c := range lk.SealCategories {
		set[string(model.ConvertSealCategory(c))] = struct{}{}
	}
	return set
}

// ---- 使用次数限制 ----
//
// 两种模式：
//   - 次数限制：License.UsageLimit > 0，每次生成章面记录一次使用，用尽后拒绝；
//   - 不限制：License.UsageLimit == 0，不校验也不计数（RecordUsage 为空操作）。
//
// 已使用次数持久化在本地 usage.json（按 LicenseId 区分），可通过环境变量
// OS_USAGE_FILE 指定记录文件路径，默认为 ~/.official-seal/usage.json。

// usageRecord 单个许可证的使用次数记录。
type usageRecord struct {
	UsedCount int64 `json:"used_count"` // 已使用次数
	UpdatedAt int64 `json:"updated_at"` // 最近一次使用时间的 Unix 时间戳
}

// usageState 使用次数持久化状态，按 LicenseId 索引。
type usageState struct {
	Records map[string]*usageRecord `json:"records"`
}

// usageMu 保护使用次数记录文件的读改写（进程内并发安全）。
var usageMu sync.Mutex

// usageFilePath 返回使用次数记录文件路径：环境变量 OS_USAGE_FILE 优先，默认 ~/.official-seal/usage.json。
func usageFilePath() (string, error) {
	if path := strings.TrimSpace(os.Getenv("OS_USAGE_FILE")); path != "" {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir for usage record: %w", err)
	}
	return filepath.Join(home, ".official-seal", "usage.json"), nil
}

// loadUsageState 读取使用次数记录，文件不存在时返回空状态。
func loadUsageState() (*usageState, error) {
	path, err := usageFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &usageState{}, nil
		}
		return nil, fmt.Errorf("read usage record %s: %w", path, err)
	}
	state := &usageState{}
	if err := json.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("parse usage record %s: %w", path, err)
	}
	return state, nil
}

// saveUsageState 原子写回使用次数记录文件（临时文件 + rename），权限 0600。
func saveUsageState(state *usageState) error {
	path, err := usageFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create usage record dir: %w", err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("write usage record %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("replace usage record %s: %w", path, err)
	}
	return nil
}

// record 返回指定许可证的记录（不存在时自动初始化）。
func (s *usageState) record(licenseId string) *usageRecord {
	if s.Records == nil {
		s.Records = make(map[string]*usageRecord)
	}
	rec, ok := s.Records[licenseId]
	if !ok || rec == nil {
		rec = &usageRecord{}
		s.Records[licenseId] = rec
	}
	return rec
}

// usageLimited 判断是否为"次数限制"模式（UsageLimit > 0）；0 表示不限制。
func (lk *License) usageLimited() bool {
	return lk != nil && lk.UsageLimit > 0
}

// usageExceeded 判断当前使用次数是否已用尽（不消耗次数）。
func (lk *License) usageExceeded() bool {
	return lk.checkUsageQuota() != nil
}

// checkUsageQuota 校验许可证使用次数是否已用尽；未限制次数时始终通过。
func (lk *License) checkUsageQuota() error {
	if !lk.usageLimited() {
		return nil
	}
	usageMu.Lock()
	defer usageMu.Unlock()
	state, err := loadUsageState()
	if err != nil {
		return err
	}
	rec := state.record(lk.LicenseId)
	if rec.UsedCount >= lk.UsageLimit {
		return fmt.Errorf("%w: used %d of %d", ErrUsageLimitExceeded, rec.UsedCount, lk.UsageLimit)
	}
	return nil
}

// recordUsage 记录一次使用（已使用次数 +1 并持久化）；未限制次数时不做任何事。
func (lk *License) recordUsage() error {
	if !lk.usageLimited() {
		return nil
	}
	usageMu.Lock()
	defer usageMu.Unlock()
	state, err := loadUsageState()
	if err != nil {
		return err
	}
	rec := state.record(lk.LicenseId)
	if rec.UsedCount >= lk.UsageLimit {
		return fmt.Errorf("%w: used %d of %d", ErrUsageLimitExceeded, rec.UsedCount, lk.UsageLimit)
	}
	rec.UsedCount++
	rec.UpdatedAt = time.Now().Unix()
	if err := saveUsageState(state); err != nil {
		Log.Error("Unable to save license usage record: %v", err)
		return err
	}
	return nil
}

// CheckUsageQuota 校验当前许可证的使用次数是否已用尽（只检查，不消耗次数）。
func CheckUsageQuota() error {
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		return ErrLicenseRequired
	}
	return lk.checkUsageQuota()
}

// RecordUsage 为当前许可证记录一次使用（已使用次数 +1 并持久化）。
// 应在章面校验通过并成功生成后调用；未限制次数时为空操作。
func RecordUsage() error {
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		return ErrLicenseRequired
	}
	return lk.recordUsage()
}

// UsageStats 返回当前许可证的使用统计：已使用次数、次数上限（0 表示不限制）。
func UsageStats() (used int64, limit int64, err error) {
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		return 0, 0, ErrLicenseRequired
	}
	limit = lk.UsageLimit
	if !lk.usageLimited() {
		return 0, limit, nil
	}
	usageMu.Lock()
	defer usageMu.Unlock()
	state, err := loadUsageState()
	if err != nil {
		return 0, limit, err
	}
	return state.record(lk.LicenseId).UsedCount, limit, nil
}
