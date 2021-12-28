package challenges

import "testing"

// https://leetcode.com/problems/best-time-to-buy-and-sell-stock/
// You want to maximize your profit by choosing a single day to buy one stock and choosing a different day in the future to sell that stock.
// Return the maximum profit you can achieve from this transaction. If you cannot achieve any profit, return 0.
func TestBestTimeToBuyAndSellStock(t *testing.T) {
	type Results struct {
		expected int
		input    []int
	}
	testCases := []Results{
		{5, []int{7, 1, 5, 3, 6, 4}},
		{0, []int{7, 6, 4, 3, 1}},
		{3, []int{7, 8, 4, 7, 1}},
		{1, []int{7, 8, 4, 3, 1}},
		{4, []int{3, 3, 5, 0, 0, 3, 1, 4}},
	}

	// maxProfit function is actual implementation
	maxProfit := func(prices []int) int {
		var profit int
		var min int = 10001
		for i := 0; i < len(prices); i++ {
			if prices[i] < min {
				min = prices[i]
			}
			if profit < prices[i]-min {
				profit = prices[i] - min
			}
		}
		return profit
	}

	for i, tc := range testCases {
		actual := maxProfit(tc.input)
		if actual != tc.expected {
			t.Errorf("Tesstcase [%d] failed, Expected [%d] but wrongly calculated to [%d]", i, tc.expected, actual)
		}
	}
}

// https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-transaction-fee
// Find the maximum profit you can achieve with one or multiple transactions considering fee per each transaction.
// Note: You may not engage in multiple transactions simultaneously (i.e., you must sell the stock before you buy again).
func TestBestTimeToBuyAndSellStockWithFee(t *testing.T) {
	type Results struct {
		expected int
		input    []int
		fee      int
	}
	testCases := []Results{
		{8, []int{1, 3, 2, 8, 4, 9}, 2},
		{6, []int{1, 3, 7, 5, 10, 3}, 3},
		{7, []int{1, 3, 6, 7, 5, 10, 3}, 2},
	}

	// maxProfitWithFee function is actual implementation
	maxProfitWithFee := func(prices []int, fee int) int {
		var profit int
		var diff int = prices[0]

		for i := 1; i < len(prices); i++ {
			diff = min(diff, prices[i]-profit)
			profit = max(profit, prices[i]-diff-fee)
		}
		return profit
	}

	for i, tc := range testCases {
		actual := maxProfitWithFee(tc.input, tc.fee)
		if actual != tc.expected {
			t.Errorf("Tesstcase [%d] failed, Expected [%d] but wrongly calculated to [%d]", i, tc.expected, actual)
		}
	}
}

// https://leetcode.com/problems/house-robber
// Each house has a certain amount of money stashed. You have to rob maximum money provided you won't alert police by robbing 2 adjacent houses.
func TestHouseRobber(t *testing.T) {
	type Results struct {
		expected int
		input    []int
	}
	testCases := []Results{
		{6, []int{3, 2, 2, 3}},
		{4, []int{1, 2, 3, 1}},
		{12, []int{2, 7, 9, 3, 1}},
		{20, []int{1, 3, 6, 7, 5, 10, 3}},
	}

	// rob function is actual implementation
	rob := func(nums []int) int {
		var output []int = make([]int, len(nums)+2)
		for i := len(nums) - 1; i >= 0; i-- {
			output[i] = max(output[i+2]+nums[i], output[i+1])
		}
		return output[0]
	}

	for i, tc := range testCases {
		actual := rob(tc.input)
		if actual != tc.expected {
			t.Errorf("Tesstcase [%d] failed, Expected [%d] but wrongly calculated to [%d]", i, tc.expected, actual)
		}
	}
}

func TestBinarySearch(t *testing.T) {
	type Results struct {
		expected int
		input    []int
		target   int
	}
	testCases := []Results{
		{4, []int{-1, 0, 3, 5, 9, 13}, 9},
		{0, []int{-1, 0, 3, 5, 9, 12, 13}, -1},
		{3, []int{-1, 0, 3, 5, 9, 12, 13}, 5},
		{-1, []int{-1, 0, 3, 5, 9, 12, 13}, 11},
		{6, []int{-1, 0, 3, 5, 9, 12, 13}, 13},
	}

	search := func(nums []int, target int) int {
		var left, pivot, right int = 0, 0, len(nums) - 1

		for left <= right {
			pivot = left + (right-left)/2
			if nums[pivot] == target {
				return pivot
			}
			if target < nums[pivot] {
				right = pivot - 1
			} else {
				left = pivot + 1
			}
		}
		return -1
	}

	for i, tc := range testCases {
		actual := search(tc.input, tc.target)
		if actual != tc.expected {
			t.Errorf("Tesstcase [%d] failed, Expected [%d] but wrongly calculated to [%d]", i, tc.expected, actual)
		}
	}
}

var min = func(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var max = func(a, b int) int {
	if a > b {
		return a
	}
	return b
}
