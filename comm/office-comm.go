package comm

func GenExcelColumn(num int) string {
	result := ""
	for num > 0 {
		num-- // 转换为从0开始的索引，因为A=1但在计算中我们从0开始
		remainder := num % 26
		result = string('A'+remainder) + result
		num /= 26
	}
	return result
}
