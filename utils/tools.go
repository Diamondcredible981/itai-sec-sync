package utils

import (
	"strconv"
	"strings"
)

// 工具函数：将逗号分隔的字符串转为整型数组
func StringToIntSlice(input string) []int {
	if input == "" {
		return []int{}
	}
	parts := strings.Split(input, ",")
	result := []int{}
	for _, part := range parts {
		if val, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			result = append(result, val)
		}
	}
	return result
}

// 工具函数：将整型数组转为逗号分隔的字符串
func IntSliceToString(input []int) string {
	parts := []string{}
	for _, val := range input {
		parts = append(parts, strconv.Itoa(val))
	}
	return strings.Join(parts, ",")
}

// 验证整数数组中的值是否都存在于给定的有效值列表中
func ValidateIntSlice(values []int, validValues []int) bool {
	validMap := make(map[int]bool)
	for _, v := range validValues {
		validMap[v] = true
	}

	for _, value := range values {
		if !validMap[value] {
			return false
		}
	}
	return true
}

// 去重整数数组
func UniqueIntSlice(slice []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range slice {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}
