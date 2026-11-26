package dto

// SealCategory 章面样式分类。
type SealCategory string

const (
	// SealCategoryCommunist 党委公章
	SealCategoryCommunist SealCategory = "communist"
	// SealCategoryEnterprise 普通公章（CompanySealCreateService.generateImage 默认类型）
	SealCategoryEnterprise SealCategory = "enterprise"
	// SealCategoryLabour 工会公章
	SealCategoryLabour SealCategory = "labour"
	// SealCategoryYouthLeague 共青团
	SealCategoryYouthLeague SealCategory = "youth_league"
	// SealCategoryCPPCC 中国人民政治协商会议章
	SealCategoryCPPCC SealCategory = "CPPCC"
	// SealCategoryOval 椭圆章
	SealCategoryOval SealCategory = "OVAL"
	// SealCategoryRectangleTimestamp 长方形时间戳章
	SealCategoryRectangleTimestamp SealCategory = "rectangle_timestamp"
	// SealCategoryCSJShusp 长三角一网通办
	SealCategoryCSJShusp SealCategory = "csj_shusp"
	// SealCategoryMoreFootCircle 多横排文字圆
	SealCategoryMoreFootCircle SealCategory = "more_foot_circle"
	// SealCategoryMoreFootRectangle 多横排文字长方形章
	SealCategoryMoreFootRectangle SealCategory = "more_foot_rectangle"
	// SealCategoryDoubleCir 双圆章
	SealCategoryDoubleCir SealCategory = "double_cir"
	// SealCategoryPublicSecurity 公检法司
	SealCategoryPublicSecurity SealCategory = "public_security"
	// SealCategoryRectangle 普通长方形章
	SealCategoryRectangle SealCategory = "rectangle"
	// SealCategorySquare 正方形章
	SealCategorySquare SealCategory = "square"
	// SealCategoryDoublelineOval 双行文字椭圆
	SealCategoryDoublelineOval SealCategory = "doubleline_oval"
	// SealCategoryDoubleOval 双椭圆章
	SealCategoryDoubleOval SealCategory = "double_oval"
	// SealCategoryArchives 档案章
	SealCategoryArchives SealCategory = "archives"
	// SealCategoryRectangleLand 规土局方章
	SealCategoryRectangleLand SealCategory = "rectangle_land"
	// SealCategoryDoublelineCir 双行文字圆章
	SealCategoryDoublelineCir SealCategory = "doubleline_cir"
	// SealCategoryMoreFootOval 多横排文字椭圆
	SealCategoryMoreFootOval SealCategory = "more_foot_oval"
	// SealCategoryPurposeSeal 用途章
	SealCategoryPurposeSeal SealCategory = "purpose_seal"
)
