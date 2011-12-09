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
  RankOrder["SGT"] = 18
  RankOrder["SGT(AAFC)"] = 18
  RankOrder["CPL"] = 19
  RankOrder["CPL(AAFC)"] = 19
  RankOrder["LAC"] = 20
  RankOrder["LAC(AAFC)"] = 20
  RankOrder["LACW"] = 20
  RankOrder["LACW(AAFC)"] = 20
  RankOrder["AC"] = 21
  RankOrder["AC(AAFC)"] = 21
  RankOrder["ACW"] = 21
  RankOrder["ACW(AAFC)"] = 21
  RankOrder["NCOCDT"] = 22
  RankOrder["ACR"] = 23
  RankOrder["ACWR"] = 23

  RankOrder["CIV"] = 25
  RankOrder["MR"] = 25
  RankOrder["MRS"] = 25
  RankOrder["MS"] = 25

  RankOrder["CUO"] = 32
  RankOrder["CWOFF"] = 33
  RankOrder["CFSGT"] = 34
  RankOrder["CSGT"] = 35
  RankOrder["CCPL"] = 36
  RankOrder["LCDT"] = 37
  RankOrder["CDT"] = 38
}
