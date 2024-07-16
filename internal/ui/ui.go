package ui

import (
	"fmt"
)

// DisplayInfo renders styled information
func DisplayInfo(info string) {
	fmt.Println(InfoStyle.Render(info))
}

// DisplayError renders a styled error message
func DisplayError(err error) {
	fmt.Println(ErrorStyle.Render(err.Error()))
}
