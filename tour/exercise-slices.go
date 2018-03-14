package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	pic := make([][]uint8, dy)
	for i := 0; i < dy; i++ {
		for k := 0; k < dx; k++ {
			pic[i] = append(pic[i], uint8((dx+dy)/2))
		}
		pic = append(pic, pic[i])
	}
	return pic
}

func main() {
	pic.Show(Pic)
}
