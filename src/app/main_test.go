package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestGetCronTabLinked(t *testing.T) {
	testScheduleData := []struct {
		x string
		XLinked string
	}{
		{x: "* * * * *", XLinked:"https://crontab.guru/#*_*_*_*_*"},
		{x: "0 * * * *", XLinked:"https://crontab.guru/#0_*_*_*_*"},
		{x: "0 * 0/3 * *", XLinked:"https://crontab.guru/#0_*_0/3_*_*"},
	}

	for _, testData := range testScheduleData {
		testedData := getCronTabLinked(testData.x)

		assert.Equal(t, testData.XLinked, testedData)
	}
}