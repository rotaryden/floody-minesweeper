package main

import "strings"

func getPrintableField(f IField) []string {
	h := f.GetHeight()
	w := f.GetWidth()
	result := make([]string, h)
	for y := 0; y < h; y++ {
		var sb strings.Builder
		for x := 0; x < w; x++ {
			pc := f.GetCell(x, y)
			switch{
			case pc.State == CellStateClosed:
				sb.WriteRune('🙫')
			case pc.HolesNumber == ThisIsHoleMarker:
				sb.WriteRune('⦿')		
			case pc.HolesNumber == 0:
				sb.WriteRune('⛶')	
			default:
				sb.WriteRune(rune(pc.HolesNumber))
			}
		}
		result[y] = sb.String()
	}
	return result
}

