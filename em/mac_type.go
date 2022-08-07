package em

type MACType int

const (
	MACTypeUndefined MACType = iota
	MACType82542
	MACType82543
	MACType82544
	MACType82540
	MACType82545
	MACType82545Rev3
	MACType82546
	MACType82546Rev3
	MACType82541
	MACType82541Rev2
	MACType82547
	MACType82547Rev2
	MACType82571
	MACType82572
	MACType82573
	MACType82574
	MACType82583
	MACType80003es2lan
	MACTypeIch8lan
	MACTypeIch9lan
	MACTypeIch10lan
	MACTypePchlan
	MACTypePch2lan
	MACTypePch_lpt
	MACTypePch_spt
	MACTypePch_cnp
	MACTypePch_tgp
	MACTypePch_adp
	MACTypePch_mtp
	MACType82575
	MACType82576
	MACType82580
	MACTypeI350
	MACTypeI354
	MACTypeI210
	MACTypeI211
	MACTypeVfadapt
	MACTypeVfadaptI350
	// List is 1-based, so subtract 1 for true count.
	NumMACType
)

// s32 e1000_set_mac_type(struct e1000_hw *hw)
func MACTypeGet(devid DeviceID) MACType {
	switch devid {
	case DEV_ID_82542:
		return MACType82542
	case DEV_ID_82543GC_FIBER:
		return MACType82543
	case DEV_ID_82543GC_COPPER:
		return MACType82543
	case DEV_ID_82544EI_COPPER:
		return MACType82544
	case DEV_ID_82544EI_FIBER:
		return MACType82544
	case DEV_ID_82544GC_COPPER:
		return MACType82544
	case DEV_ID_82544GC_LOM:
		return MACType82544
	case DEV_ID_82540EM:
		return MACType82540
	case DEV_ID_82540EM_LOM:
		return MACType82540
	case DEV_ID_82540EP:
		return MACType82540
	case DEV_ID_82540EP_LOM:
		return MACType82540
	case DEV_ID_82540EP_LP:
		return MACType82540
	case DEV_ID_82545EM_COPPER:
		return MACType82545
	case DEV_ID_82545EM_FIBER:
		return MACType82545
	case DEV_ID_82545GM_COPPER:
		return MACType82545Rev3
	case DEV_ID_82545GM_FIBER:
		return MACType82545Rev3
	case DEV_ID_82545GM_SERDES:
		return MACType82545Rev3
	case DEV_ID_82546EB_COPPER:
		return MACType82546
	case DEV_ID_82546EB_FIBER:
		return MACType82546
	case DEV_ID_82546EB_QUAD_COPPER:
		return MACType82546
	case DEV_ID_82546GB_COPPER:
		return MACType82546Rev3
	case DEV_ID_82546GB_FIBER:
		return MACType82546Rev3
	case DEV_ID_82546GB_SERDES:
		return MACType82546Rev3
	case DEV_ID_82546GB_PCIE:
		return MACType82546Rev3
	case DEV_ID_82546GB_QUAD_COPPER:
		return MACType82546Rev3
	case DEV_ID_82546GB_QUAD_COPPER_KSP3:
		return MACType82546Rev3
	case DEV_ID_82541EI:
		return MACType82541
	case DEV_ID_82541EI_MOBILE:
		return MACType82541
	case DEV_ID_82541ER_LOM:
		return MACType82541
	case DEV_ID_82541ER:
		return MACType82541Rev2
	case DEV_ID_82541GI:
		return MACType82541Rev2
	case DEV_ID_82541GI_LF:
		return MACType82541Rev2
	case DEV_ID_82541GI_MOBILE:
		return MACType82541Rev2
	case DEV_ID_82547EI:
		return MACType82547
	case DEV_ID_82547EI_MOBILE:
		return MACType82547
	case DEV_ID_82547GI:
		return MACType82547Rev2
	case DEV_ID_82571EB_COPPER:
		return MACType82571
	case DEV_ID_82571EB_FIBER:
		return MACType82571
	case DEV_ID_82571EB_SERDES:
		return MACType82571
	case DEV_ID_82571EB_SERDES_DUAL:
		return MACType82571
	case DEV_ID_82571EB_SERDES_QUAD:
		return MACType82571
	case DEV_ID_82571EB_QUAD_COPPER:
		return MACType82571
	case DEV_ID_82571PT_QUAD_COPPER:
		return MACType82571
	case DEV_ID_82571EB_QUAD_FIBER:
		return MACType82571
	case DEV_ID_82571EB_QUAD_COPPER_LP:
		return MACType82571
	case DEV_ID_82572EI:
		return MACType82572
	case DEV_ID_82572EI_COPPER:
		return MACType82572
	case DEV_ID_82572EI_FIBER:
		return MACType82572
	case DEV_ID_82572EI_SERDES:
		return MACType82572
	case DEV_ID_82573E:
		return MACType82573
	case DEV_ID_82573E_IAMT:
		return MACType82573
	case DEV_ID_82573L:
		return MACType82573
	case DEV_ID_82574L:
		return MACType82574
	case DEV_ID_82574LA:
		return MACType82574
	case DEV_ID_82583V:
		return MACType82583
	case DEV_ID_80003ES2LAN_COPPER_DPT:
		return MACType80003es2lan
	case DEV_ID_80003ES2LAN_SERDES_DPT:
		return MACType80003es2lan
	case DEV_ID_80003ES2LAN_COPPER_SPT:
		return MACType80003es2lan
	case DEV_ID_80003ES2LAN_SERDES_SPT:
		return MACType80003es2lan
	case DEV_ID_ICH8_IFE:
		return MACTypeIch8lan
	case DEV_ID_ICH8_IFE_GT:
		return MACTypeIch8lan
	case DEV_ID_ICH8_IFE_G:
		return MACTypeIch8lan
	case DEV_ID_ICH8_IGP_M:
		return MACTypeIch8lan
	case DEV_ID_ICH8_IGP_M_AMT:
		return MACTypeIch8lan
	case DEV_ID_ICH8_IGP_AMT:
		return MACTypeIch8lan
	case DEV_ID_ICH8_IGP_C:
		return MACTypeIch8lan
	case DEV_ID_ICH8_82567V_3:
		return MACTypeIch8lan
	case DEV_ID_ICH9_IFE:
		return MACTypeIch9lan
	case DEV_ID_ICH9_IFE_GT:
		return MACTypeIch9lan
	case DEV_ID_ICH9_IFE_G:
		return MACTypeIch9lan
	case DEV_ID_ICH9_IGP_M:
		return MACTypeIch9lan
	case DEV_ID_ICH9_IGP_M_AMT:
		return MACTypeIch9lan
	case DEV_ID_ICH9_IGP_M_V:
		return MACTypeIch9lan
	case DEV_ID_ICH9_IGP_AMT:
		return MACTypeIch9lan
	case DEV_ID_ICH9_BM:
		return MACTypeIch9lan
	case DEV_ID_ICH9_IGP_C:
		return MACTypeIch9lan
	case DEV_ID_ICH10_R_BM_LM:
		return MACTypeIch9lan
	case DEV_ID_ICH10_R_BM_LF:
		return MACTypeIch9lan
	case DEV_ID_ICH10_R_BM_V:
		return MACTypeIch9lan
	case DEV_ID_ICH10_D_BM_LM:
		return MACTypeIch10lan
	case DEV_ID_ICH10_D_BM_LF:
		return MACTypeIch10lan
	case DEV_ID_ICH10_D_BM_V:
		return MACTypeIch10lan
	case DEV_ID_PCH_D_HV_DM:
		return MACTypePchlan
	case DEV_ID_PCH_D_HV_DC:
		return MACTypePchlan
	case DEV_ID_PCH_M_HV_LM:
		return MACTypePchlan
	case DEV_ID_PCH_M_HV_LC:
		return MACTypePchlan
	case DEV_ID_PCH2_LV_LM:
		return MACTypePch2lan
	case DEV_ID_PCH2_LV_V:
		return MACTypePch2lan
	case DEV_ID_PCH_LPT_I217_LM:
		return MACTypePch_lpt
	case DEV_ID_PCH_LPT_I217_V:
		return MACTypePch_lpt
	case DEV_ID_PCH_LPTLP_I218_LM:
		return MACTypePch_lpt
	case DEV_ID_PCH_LPTLP_I218_V:
		return MACTypePch_lpt
	case DEV_ID_PCH_I218_LM2:
		return MACTypePch_lpt
	case DEV_ID_PCH_I218_V2:
		return MACTypePch_lpt
	case DEV_ID_PCH_I218_LM3:
		return MACTypePch_lpt
	case DEV_ID_PCH_I218_V3:
		return MACTypePch_lpt
	case DEV_ID_PCH_SPT_I219_LM:
		return MACTypePch_spt
	case DEV_ID_PCH_SPT_I219_V:
		return MACTypePch_spt
	case DEV_ID_PCH_SPT_I219_LM2:
		return MACTypePch_spt
	case DEV_ID_PCH_SPT_I219_V2:
		return MACTypePch_spt
	case DEV_ID_PCH_LBG_I219_LM3:
		return MACTypePch_spt
	case DEV_ID_PCH_SPT_I219_LM4:
		return MACTypePch_spt
	case DEV_ID_PCH_SPT_I219_V4:
		return MACTypePch_spt
	case DEV_ID_PCH_SPT_I219_LM5:
		return MACTypePch_spt
	case DEV_ID_PCH_SPT_I219_V5:
		return MACTypePch_spt
	case DEV_ID_PCH_CMP_I219_LM12:
		return MACTypePch_spt
	case DEV_ID_PCH_CMP_I219_V12:
		return MACTypePch_spt

	case DEV_ID_PCH_CNP_I219_LM6:
		return MACTypePch_cnp
	case DEV_ID_PCH_CNP_I219_V6:
		return MACTypePch_cnp
	case DEV_ID_PCH_CNP_I219_LM7:
		return MACTypePch_cnp
	case DEV_ID_PCH_CNP_I219_V7:
		return MACTypePch_cnp
	case DEV_ID_PCH_ICP_I219_LM8:
		return MACTypePch_cnp
	case DEV_ID_PCH_ICP_I219_V8:
		return MACTypePch_cnp
	case DEV_ID_PCH_ICP_I219_LM9:
		return MACTypePch_cnp
	case DEV_ID_PCH_ICP_I219_V9:
		return MACTypePch_cnp
	case DEV_ID_PCH_CMP_I219_LM10:
		return MACTypePch_cnp
	case DEV_ID_PCH_CMP_I219_V10:
		return MACTypePch_cnp
	case DEV_ID_PCH_CMP_I219_LM11:
		return MACTypePch_cnp
	case DEV_ID_PCH_CMP_I219_V11:
		return MACTypePch_cnp

	case DEV_ID_PCH_TGP_I219_LM13:
		return MACTypePch_tgp
	case DEV_ID_PCH_TGP_I219_V13:
		return MACTypePch_tgp
	case DEV_ID_PCH_TGP_I219_LM14:
		return MACTypePch_tgp
	case DEV_ID_PCH_TGP_I219_V14:
		return MACTypePch_tgp
	case DEV_ID_PCH_TGP_I219_LM15:
		return MACTypePch_tgp
	case DEV_ID_PCH_TGP_I219_V15:
		return MACTypePch_tgp

	case DEV_ID_PCH_ADL_I219_LM16:
		return MACTypePch_adp
	case DEV_ID_PCH_ADL_I219_V16:
		return MACTypePch_adp
	case DEV_ID_PCH_ADL_I219_LM17:
		return MACTypePch_adp
	case DEV_ID_PCH_ADL_I219_V17:
		return MACTypePch_adp

	case DEV_ID_PCH_MTP_I219_LM18:
		return MACTypePch_mtp
	case DEV_ID_PCH_MTP_I219_V18:
		return MACTypePch_mtp
	case DEV_ID_PCH_MTP_I219_LM19:
		return MACTypePch_mtp
	case DEV_ID_PCH_MTP_I219_V19:
		return MACTypePch_mtp

	case DEV_ID_82575EB_COPPER:
		return MACType82575
	case DEV_ID_82575EB_FIBER_SERDES:
		return MACType82575
	case DEV_ID_82575GB_QUAD_COPPER:
		return MACType82575

	case DEV_ID_82576:
		return MACType82576
	case DEV_ID_82576_FIBER:
		return MACType82576
	case DEV_ID_82576_SERDES:
		return MACType82576
	case DEV_ID_82576_QUAD_COPPER:
		return MACType82576
	case DEV_ID_82576_QUAD_COPPER_ET2:
		return MACType82576
	case DEV_ID_82576_NS:
		return MACType82576
	case DEV_ID_82576_NS_SERDES:
		return MACType82576
	case DEV_ID_82576_SERDES_QUAD:
		return MACType82576

	case DEV_ID_82580_COPPER:
		return MACType82580
	case DEV_ID_82580_FIBER:
		return MACType82580
	case DEV_ID_82580_SERDES:
		return MACType82580
	case DEV_ID_82580_SGMII:
		return MACType82580
	case DEV_ID_82580_COPPER_DUAL:
		return MACType82580
	case DEV_ID_82580_QUAD_FIBER:
		return MACType82580
	case DEV_ID_DH89XXCC_SGMII:
		return MACType82580
	case DEV_ID_DH89XXCC_SERDES:
		return MACType82580
	case DEV_ID_DH89XXCC_BACKPLANE:
		return MACType82580
	case DEV_ID_DH89XXCC_SFP:
		return MACType82580

	case DEV_ID_I350_COPPER:
		return MACTypeI350
	case DEV_ID_I350_FIBER:
		return MACTypeI350
	case DEV_ID_I350_SERDES:
		return MACTypeI350
	case DEV_ID_I350_SGMII:
		return MACTypeI350
	case DEV_ID_I350_DA4:
		return MACTypeI350

	case DEV_ID_I210_COPPER_FLASHLESS:
		return MACTypeI210
	case DEV_ID_I210_SERDES_FLASHLESS:
		return MACTypeI210
	case DEV_ID_I210_SGMII_FLASHLESS:
		return MACTypeI210
	case DEV_ID_I210_COPPER:
		return MACTypeI210
	case DEV_ID_I210_COPPER_OEM1:
		return MACTypeI210
	case DEV_ID_I210_COPPER_IT:
		return MACTypeI210
	case DEV_ID_I210_FIBER:
		return MACTypeI210
	case DEV_ID_I210_SERDES:
		return MACTypeI210
	case DEV_ID_I210_SGMII:
		return MACTypeI210

	case DEV_ID_I211_COPPER:
		return MACTypeI211

	case DEV_ID_82576_VF:
		return MACTypeVfadapt
	case DEV_ID_82576_VF_HV:
		return MACTypeVfadapt

	case DEV_ID_I350_VF:
		return MACTypeVfadaptI350
	case DEV_ID_I350_VF_HV:
		return MACTypeVfadaptI350

	case DEV_ID_I354_BACKPLANE_1GBPS:
		return MACTypeI354
	case DEV_ID_I354_SGMII:
		return MACTypeI354
	case DEV_ID_I354_BACKPLANE_2_5GBPS:
		return MACTypeI354
	default:
		return MACTypeUndefined
	}
}

// int em_set_num_queues(if_ctx_t ctx)
func (t MACType) NumQueue() int {
	switch t {
	case MACType82576:
		return 8
	case MACType82580:
		return 8
	case MACTypeI350:
		return 8
	case MACTypeI354:
		return 8
	case MACTypeI210:
		return 4
	case MACType82575:
		return 4
	case MACTypeI211:
		return 2
	case MACType82574:
		return 2
	default:
		return 1
	}
}
