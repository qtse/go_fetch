package movo

var (
    RankOrder map[string]int = make(map[string]int)
    )

func init() {
  RankOrder["ACM"] = 1
  RankOrder["AIRMSHL"] = 2
  RankOrder["AVM"] = 3
  RankOrder["WOFF-AF"] = 4
  RankOrder["AIRCDRE"] = 5

  RankOrder["GPCAPT"] = 6
  RankOrder["GPCAPT(AAFC)"] = 6
  RankOrder["WGCDR"] = 7
  RankOrder["WGCDR(AAFC)"] = 7
  RankOrder["SQNLDR"] = 8
  RankOrder["SQNLDR(AAFC)"] = 8
  RankOrder["FLTLT"] = 9
  RankOrder["FLTLT(AAFC)"] = 9
  RankOrder["FLGOFF"] = 10
  RankOrder["FLGOFF(AAFC)"] = 10
  RankOrder["PLTOFF"] = 11
  RankOrder["PLTOFF(AAFC)"] = 11
  RankOrder["OFFCDT"] = 12

  RankOrder["WOFF"] = 16
  RankOrder["WOFF(AAFC)"] = 16
  RankOrder["FSGT"] = 17
  RankOrder["FSGT(AAFC)"] = 17
  RankOrder["SGT"] = 19
  RankOrder["SGT(AAFC)"] = 19
  RankOrder["CPL"] = 20
  RankOrder["CPL(AAFC)"] = 20
  RankOrder["LAC"] = 22
  RankOrder["LAC(AAFC)"] = 22
  RankOrder["LACW"] = 22
  RankOrder["LACW(AAFC)"] = 22
  RankOrder["AC"] = 23
  RankOrder["AC(AAFC)"] = 23
  RankOrder["ACW"] = 23
  RankOrder["ACW(AAFC)"] = 23
  RankOrder["NCOCDT"] = 24
  RankOrder["ACR"] = 25
  RankOrder["ACWR"] = 25

  RankOrder["CIV"] = 26
  RankOrder["MR"] = 26
  RankOrder["MRS"] = 26
  RankOrder["MS"] = 26

  RankOrder["CUO"] = 32
  RankOrder["CWOFF"] = 33
  RankOrder["CFSGT"] = 34
  RankOrder["CSGT"] = 36
  RankOrder["CCPL"] = 37
  RankOrder["LCDT"] = 39
  RankOrder["CDT"] = 40

  RankOrder["GEN"] = 1
  RankOrder["LTGEN"] = 2
  RankOrder["MAJGEN"] = 3
  RankOrder["RSM-A"] = 4
  RankOrder["BRIG"] = 5

  RankOrder["COL"] = 6
  RankOrder["COL(AAC)"] = 6
  RankOrder["LTCOL"] = 7
  RankOrder["LTCOL(AAC)"] = 7
  RankOrder["MAJ"] = 8
  RankOrder["MAJ(AAC)"] = 8
///  RankOrder["CAPT"] = 9
  RankOrder["CAPT(AAC)"] = 9
  RankOrder["LT"] = 10
  RankOrder["LT(AAC)"] = 10
  RankOrder["2LT"] = 11
  RankOrder["2LT(AAC)"] = 11

  RankOrder["WO1"] = 16
  RankOrder["WO1(AAC)"] = 16
  RankOrder["WO2"] = 17
  RankOrder["WO2(AAC)"] = 16
  RankOrder["SSGT"] = 18
  RankOrder["SSGT(AAC)"] = 18
///  RankOrder["SGT"] = 19
  RankOrder["SGT(AAC)"] = 19
///  RankOrder["CPL"] = 20
  RankOrder["CPL(AAC)"] = 20
  RankOrder["LCPL"] = 21
  RankOrder["LCPL(AAC)"] = 21
  RankOrder["PTE"] = 23
  RankOrder["UA(AAC)"] = 23
  RankOrder["OCDT"] = 24
  RankOrder["RCT"] = 25

  RankOrder["NATCUO"] = 32
  RankOrder["RCUO"] = 32
///  RankOrder["CUO"] = 32
  RankOrder["NATCDTRSM"] = 33
  RankOrder["CDTWO1"] = 33
  RankOrder["CDTWO2"] = 34
  RankOrder["CDTSSGT"] = 35
  RankOrder["CDTSGT"] = 36
  RankOrder["CDTCPL"] = 37
  RankOrder["CDTLCPL"] = 38
///  RankOrder["CDT"] = 39
  RankOrder["CDTREC"] = 40

  RankOrder["ADML"] = 1
  RankOrder["VADM"] = 2
  RankOrder["RADM"] = 3
  RankOrder["WO-N"] = 4
  RankOrder["CDRE"] = 5

  RankOrder["CAPT"] = 6
  RankOrder["CMDR"] = 7
  RankOrder["LCDR"] = 8
  RankOrder["LEUT"] = 9
  RankOrder["SBLT"] = 10
  RankOrder["MIDN"] = 12

  RankOrder["WO"] = 16
  RankOrder["CPO"] = 17
  RankOrder["PO"] = 19
  RankOrder["LS"] = 20
  RankOrder["AB"] = 21
  RankOrder["SMN"] = 23

  RankOrder["CDTMIDN"] = 29
///  Could also be (But it's confusing):
///  RankOrder["MIDN"] = 29
  RankOrder["CDTWO"] = 33
  RankOrder["CDTCPO"] = 34
  RankOrder["CDTPO"] = 36
  RankOrder["CDTLS"] = 37
  RankOrder["CDTAB"] = 38
  RankOrder["CDTSMN"] = 39
  RankOrder["CDTRCT"] = 40
}
