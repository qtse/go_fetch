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

  RankOrder["WOFF"] = 13
  RankOrder["WOFF(AAFC)"] = 13
  RankOrder["FSGT"] = 14
  RankOrder["FSGT(AAFC)"] = 14
  RankOrder["SGT"] = 15
  RankOrder["SGT(AAFC)"] = 15
  RankOrder["CPL"] = 16
  RankOrder["CPL(AAFC)"] = 16
  RankOrder["LAC"] = 17
  RankOrder["LAC(AAFC)"] = 17
  RankOrder["LACW"] = 17
  RankOrder["LACW(AAFC)"] = 17
  RankOrder["AC"] = 18
  RankOrder["AC(AAFC)"] = 18
  RankOrder["ACW"] = 18
  RankOrder["ACW(AAFC)"] = 18
  RankOrder["NCOCDT"] = 19
  RankOrder["ACR"] = 20
  RankOrder["ACWR"] = 20

  RankOrder["CUO"] = 21
  RankOrder["CWOFF"] = 22
  RankOrder["CFSGT"] = 23
  RankOrder["CSGT"] = 24
  RankOrder["CCPL"] = 25
  RankOrder["LCDT"] = 26
  RankOrder["CDT"] = 27
}
