package util

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

const (
	// RED
	RED = "\033[31m"
	// GREEN
	GREEN = "\033[32m"
	// YELLOW
	YELLOW = "\033[33m"
	// BLUE
	BLUE = "\033[34m"
	// FUCHSIA
	FUCHSIA = "\033[35m"
	// CYAN
	CYAN = "\033[36m"
	// WHITE
	WHITE = "\033[37m"
	// RESET
	RESET = "\033[0m"
)

// IsInteger Determine whether the string is an integer
func IsInteger(input string) bool {
	_, err := strconv.Atoi(input)
	return err == nil
}

// RandString Random string
func RandString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// VerifyEmailFormat E-mail verification
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func getChar(str string) string {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()
	fmt.Print(str)
	char, _, _ := keyboard.GetKey()
	fmt.Printf("%c\n", char)
	if char == 0 {
		return ""
	}
	return string(char)
}

// LoopInput Cycle input selection, or return directly to exit
func LoopInput(tip string, choices interface{}, singleRowPrint bool) int {
	reflectValue := reflect.ValueOf(choices)
	if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
		fmt.Println("only support slice or array type!")
		return -1
	}
	length := reflectValue.Len()
	if reflectValue.Type().String() == "[]string" {
		if singleRowPrint {
			for i := 0; i < length; i++ {
				fmt.Printf("%d.%s\n\n", i+1, reflectValue.Index(i).Interface())
			}
		} else {
			for i := 0; i < length; i++ {
				if i%2 == 0 {
					fmt.Printf("%d.%-15s\t", i+1, reflectValue.Index(i).Interface())
				} else {
					fmt.Printf("%d.%-15s\n\n", i+1, reflectValue.Index(i).Interface())
				}
			}
		}
	}
	for {
		inputString := ""
		if length < 10 {
			inputString = getChar(tip)
		} else {
			fmt.Print(tip)
			_, _ = fmt.Scanln(&inputString)
		}
		if inputString == "" {
			return -1
		} else if !IsInteger(inputString) {
			fmt.Println("Input is wrong, please re-enter")
			continue
		}
		number, _ := strconv.Atoi(inputString)
		if number <= length && number > 0 {
			return number
		}
		fmt.Println("The input number is out of bounds, please re-enter")
	}
}

// Input Read terminal user input
func Input(tip string, defaultValue string) string {
	input := ""
	fmt.Print(tip)
	_, _ = fmt.Scanln(&input)
	if input == "" && defaultValue != "" {
		input = defaultValue
	}
	return input
}

// Red
func Red(str string) string {
	return RED + str + RESET
}

// Green
func Green(str string) string {
	return GREEN + str + RESET
}

// Yellow
func Yellow(str string) string {
	return YELLOW + str + RESET
}

// Blue
func Blue(str string) string {
	return BLUE + str + RESET
}

// Fuchsia
func Fuchsia(str string) string {
	return FUCHSIA + str + RESET
}

// Cyan
func Cyan(str string) string {
	return CYAN + str + RESET
}

// White
func White(str string) string {
	return WHITE + str + RESET
}
