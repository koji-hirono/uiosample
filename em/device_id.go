package em

type DeviceID uint16

const (
	DEV_ID_82542                    DeviceID = 0x1000
	DEV_ID_82543GC_FIBER                     = 0x1001
	DEV_ID_82543GC_COPPER                    = 0x1004
	DEV_ID_82544EI_COPPER                    = 0x1008
	DEV_ID_82544EI_FIBER                     = 0x1009
	DEV_ID_82544GC_COPPER                    = 0x100C
	DEV_ID_82544GC_LOM                       = 0x100D
	DEV_ID_82540EM                           = 0x100E
	DEV_ID_82540EM_LOM                       = 0x1015
	DEV_ID_82540EP_LOM                       = 0x1016
	DEV_ID_82540EP                           = 0x1017
	DEV_ID_82540EP_LP                        = 0x101E
	DEV_ID_82545EM_COPPER                    = 0x100F
	DEV_ID_82545EM_FIBER                     = 0x1011
	DEV_ID_82545GM_COPPER                    = 0x1026
	DEV_ID_82545GM_FIBER                     = 0x1027
	DEV_ID_82545GM_SERDES                    = 0x1028
	DEV_ID_82546EB_COPPER                    = 0x1010
	DEV_ID_82546EB_FIBER                     = 0x1012
	DEV_ID_82546EB_QUAD_COPPER               = 0x101D
	DEV_ID_82546GB_COPPER                    = 0x1079
	DEV_ID_82546GB_FIBER                     = 0x107A
	DEV_ID_82546GB_SERDES                    = 0x107B
	DEV_ID_82546GB_PCIE                      = 0x108A
	DEV_ID_82546GB_QUAD_COPPER               = 0x1099
	DEV_ID_82546GB_QUAD_COPPER_KSP3          = 0x10B5
	DEV_ID_82541EI                           = 0x1013
	DEV_ID_82541EI_MOBILE                    = 0x1018
	DEV_ID_82541ER_LOM                       = 0x1014
	DEV_ID_82541ER                           = 0x1078
	DEV_ID_82541GI                           = 0x1076
	DEV_ID_82541GI_LF                        = 0x107C
	DEV_ID_82541GI_MOBILE                    = 0x1077
	DEV_ID_82547EI                           = 0x1019
	DEV_ID_82547EI_MOBILE                    = 0x101A
	DEV_ID_82547GI                           = 0x1075
	DEV_ID_82571EB_COPPER                    = 0x105E
	DEV_ID_82571EB_FIBER                     = 0x105F
	DEV_ID_82571EB_SERDES                    = 0x1060
	DEV_ID_82571EB_SERDES_DUAL               = 0x10D9
	DEV_ID_82571EB_SERDES_QUAD               = 0x10DA
	DEV_ID_82571EB_QUAD_COPPER               = 0x10A4
	DEV_ID_82571PT_QUAD_COPPER               = 0x10D5
	DEV_ID_82571EB_QUAD_FIBER                = 0x10A5
	DEV_ID_82571EB_QUAD_COPPER_LP            = 0x10BC
	DEV_ID_82572EI_COPPER                    = 0x107D
	DEV_ID_82572EI_FIBER                     = 0x107E
	DEV_ID_82572EI_SERDES                    = 0x107F
	DEV_ID_82572EI                           = 0x10B9
	DEV_ID_82573E                            = 0x108B
	DEV_ID_82573E_IAMT                       = 0x108C
	DEV_ID_82573L                            = 0x109A
	DEV_ID_82574L                            = 0x10D3
	DEV_ID_82574LA                           = 0x10F6
	DEV_ID_82583V                            = 0x150C
	DEV_ID_80003ES2LAN_COPPER_DPT            = 0x1096
	DEV_ID_80003ES2LAN_SERDES_DPT            = 0x1098
	DEV_ID_80003ES2LAN_COPPER_SPT            = 0x10BA
	DEV_ID_80003ES2LAN_SERDES_SPT            = 0x10BB
	DEV_ID_ICH8_82567V_3                     = 0x1501
	DEV_ID_ICH8_IGP_M_AMT                    = 0x1049
	DEV_ID_ICH8_IGP_AMT                      = 0x104A
	DEV_ID_ICH8_IGP_C                        = 0x104B
	DEV_ID_ICH8_IFE                          = 0x104C
	DEV_ID_ICH8_IFE_GT                       = 0x10C4
	DEV_ID_ICH8_IFE_G                        = 0x10C5
	DEV_ID_ICH8_IGP_M                        = 0x104D
	DEV_ID_ICH9_IGP_M                        = 0x10BF
	DEV_ID_ICH9_IGP_M_AMT                    = 0x10F5
	DEV_ID_ICH9_IGP_M_V                      = 0x10CB
	DEV_ID_ICH9_IGP_AMT                      = 0x10BD
	DEV_ID_ICH9_BM                           = 0x10E5
	DEV_ID_ICH9_IGP_C                        = 0x294C
	DEV_ID_ICH9_IFE                          = 0x10C0
	DEV_ID_ICH9_IFE_GT                       = 0x10C3
	DEV_ID_ICH9_IFE_G                        = 0x10C2
	DEV_ID_ICH10_R_BM_LM                     = 0x10CC
	DEV_ID_ICH10_R_BM_LF                     = 0x10CD
	DEV_ID_ICH10_R_BM_V                      = 0x10CE
	DEV_ID_ICH10_D_BM_LM                     = 0x10DE
	DEV_ID_ICH10_D_BM_LF                     = 0x10DF
	DEV_ID_ICH10_D_BM_V                      = 0x1525
	DEV_ID_PCH_M_HV_LM                       = 0x10EA
	DEV_ID_PCH_M_HV_LC                       = 0x10EB
	DEV_ID_PCH_D_HV_DM                       = 0x10EF
	DEV_ID_PCH_D_HV_DC                       = 0x10F0
	DEV_ID_PCH2_LV_LM                        = 0x1502
	DEV_ID_PCH2_LV_V                         = 0x1503
	DEV_ID_PCH_LPT_I217_LM                   = 0x153A
	DEV_ID_PCH_LPT_I217_V                    = 0x153B
	DEV_ID_PCH_LPTLP_I218_LM                 = 0x155A
	DEV_ID_PCH_LPTLP_I218_V                  = 0x1559
	DEV_ID_PCH_I218_LM2                      = 0x15A0
	DEV_ID_PCH_I218_V2                       = 0x15A1
	DEV_ID_PCH_I218_LM3                      = 0x15A2 // Wildcat Point PCH
	DEV_ID_PCH_I218_V3                       = 0x15A3 // Wildcat Point PCH
	DEV_ID_PCH_SPT_I219_LM                   = 0x156F // Sunrise Point PCH
	DEV_ID_PCH_SPT_I219_V                    = 0x1570 // Sunrise Point PCH
	DEV_ID_PCH_SPT_I219_LM2                  = 0x15B7 // Sunrise Point-H PCH
	DEV_ID_PCH_SPT_I219_V2                   = 0x15B8 // Sunrise Point-H PCH
	DEV_ID_PCH_LBG_I219_LM3                  = 0x15B9 // LEWISBURG PCH
	DEV_ID_PCH_SPT_I219_LM4                  = 0x15D7
	DEV_ID_PCH_SPT_I219_V4                   = 0x15D8
	DEV_ID_PCH_SPT_I219_LM5                  = 0x15E3
	DEV_ID_PCH_SPT_I219_V5                   = 0x15D6
	DEV_ID_PCH_CNP_I219_LM6                  = 0x15BD
	DEV_ID_PCH_CNP_I219_V6                   = 0x15BE
	DEV_ID_PCH_CNP_I219_LM7                  = 0x15BB
	DEV_ID_PCH_CNP_I219_V7                   = 0x15BC
	DEV_ID_PCH_ICP_I219_LM8                  = 0x15DF
	DEV_ID_PCH_ICP_I219_V8                   = 0x15E0
	DEV_ID_PCH_ICP_I219_LM9                  = 0x15E1
	DEV_ID_PCH_ICP_I219_V9                   = 0x15E2
	DEV_ID_PCH_CMP_I219_LM10                 = 0x0D4E
	DEV_ID_PCH_CMP_I219_V10                  = 0x0D4F
	DEV_ID_PCH_CMP_I219_LM11                 = 0x0D4C
	DEV_ID_PCH_CMP_I219_V11                  = 0x0D4D
	DEV_ID_PCH_CMP_I219_LM12                 = 0x0D53
	DEV_ID_PCH_CMP_I219_V12                  = 0x0D55
	DEV_ID_PCH_TGP_I219_LM13                 = 0x15FB
	DEV_ID_PCH_TGP_I219_V13                  = 0x15FC
	DEV_ID_PCH_TGP_I219_LM14                 = 0x15F9
	DEV_ID_PCH_TGP_I219_V14                  = 0x15FA
	DEV_ID_PCH_TGP_I219_LM15                 = 0x15F4
	DEV_ID_PCH_TGP_I219_V15                  = 0x15F5
	DEV_ID_PCH_ADL_I219_LM16                 = 0x1A1E
	DEV_ID_PCH_ADL_I219_V16                  = 0x1A1F
	DEV_ID_PCH_ADL_I219_LM17                 = 0x1A1C
	DEV_ID_PCH_ADL_I219_V17                  = 0x1A1D
	DEV_ID_PCH_MTP_I219_LM18                 = 0x550A
	DEV_ID_PCH_MTP_I219_V18                  = 0x550B
	DEV_ID_PCH_MTP_I219_LM19                 = 0x550C
	DEV_ID_PCH_MTP_I219_V19                  = 0x550D
	DEV_ID_82576                             = 0x10C9
	DEV_ID_82576_FIBER                       = 0x10E6
	DEV_ID_82576_SERDES                      = 0x10E7
	DEV_ID_82576_QUAD_COPPER                 = 0x10E8
	DEV_ID_82576_QUAD_COPPER_ET2             = 0x1526
	DEV_ID_82576_NS                          = 0x150A
	DEV_ID_82576_NS_SERDES                   = 0x1518
	DEV_ID_82576_SERDES_QUAD                 = 0x150D
	DEV_ID_82576_VF                          = 0x10CA
	DEV_ID_82576_VF_HV                       = 0x152D
	DEV_ID_I350_VF                           = 0x1520
	DEV_ID_I350_VF_HV                        = 0x152F
	DEV_ID_82575EB_COPPER                    = 0x10A7
	DEV_ID_82575EB_FIBER_SERDES              = 0x10A9
	DEV_ID_82575GB_QUAD_COPPER               = 0x10D6
	DEV_ID_82580_COPPER                      = 0x150E
	DEV_ID_82580_FIBER                       = 0x150F
	DEV_ID_82580_SERDES                      = 0x1510
	DEV_ID_82580_SGMII                       = 0x1511
	DEV_ID_82580_COPPER_DUAL                 = 0x1516
	DEV_ID_82580_QUAD_FIBER                  = 0x1527
	DEV_ID_I350_COPPER                       = 0x1521
	DEV_ID_I350_FIBER                        = 0x1522
	DEV_ID_I350_SERDES                       = 0x1523
	DEV_ID_I350_SGMII                        = 0x1524
	DEV_ID_I350_DA4                          = 0x1546
	DEV_ID_I210_COPPER                       = 0x1533
	DEV_ID_I210_COPPER_OEM1                  = 0x1534
	DEV_ID_I210_COPPER_IT                    = 0x1535
	DEV_ID_I210_FIBER                        = 0x1536
	DEV_ID_I210_SERDES                       = 0x1537
	DEV_ID_I210_SGMII                        = 0x1538
	DEV_ID_I210_COPPER_FLASHLESS             = 0x157B
	DEV_ID_I210_SERDES_FLASHLESS             = 0x157C
	DEV_ID_I210_SGMII_FLASHLESS              = 0x15F6
	DEV_ID_I211_COPPER                       = 0x1539
	DEV_ID_I354_BACKPLANE_1GBPS              = 0x1F40
	DEV_ID_I354_SGMII                        = 0x1F41
	DEV_ID_I354_BACKPLANE_2_5GBPS            = 0x1F45
	DEV_ID_DH89XXCC_SGMII                    = 0x0438
	DEV_ID_DH89XXCC_SERDES                   = 0x043A
	DEV_ID_DH89XXCC_BACKPLANE                = 0x043C
	DEV_ID_DH89XXCC_SFP                      = 0x0440
)

func (id DeviceID) IsICH8() bool {
	switch id {
	case DEV_ID_PCH2_LV_LM:
	case DEV_ID_PCH_LPT_I217_LM:
	case DEV_ID_PCH_LPT_I217_V:
	case DEV_ID_PCH_LPTLP_I218_LM:
	case DEV_ID_PCH_LPTLP_I218_V:
	case DEV_ID_PCH_I218_V2:
	case DEV_ID_PCH_I218_LM2:
	case DEV_ID_PCH_I218_V3:
	case DEV_ID_PCH_I218_LM3:
	case DEV_ID_PCH_SPT_I219_LM:
	case DEV_ID_PCH_SPT_I219_V:
	case DEV_ID_PCH_SPT_I219_LM2:
	case DEV_ID_PCH_SPT_I219_V2:
	case DEV_ID_PCH_LBG_I219_LM3:
	case DEV_ID_PCH_SPT_I219_LM4:
	case DEV_ID_PCH_SPT_I219_V4:
	case DEV_ID_PCH_SPT_I219_LM5:
	case DEV_ID_PCH_SPT_I219_V5:
	case DEV_ID_PCH_CNP_I219_LM6:
	case DEV_ID_PCH_CNP_I219_V6:
	case DEV_ID_PCH_CNP_I219_LM7:
	case DEV_ID_PCH_CNP_I219_V7:
	default:
		return false
	}
	return true
}
