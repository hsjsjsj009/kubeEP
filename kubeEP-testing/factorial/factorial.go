package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

func factorial(start *big.Int) *big.Int {
	one := big.NewInt(1)
	res := big.NewInt(1)
	for x := big.NewInt(2); x.CmpAbs(start) < 1; x = x.Add(x, one) {
		res = res.Mul(res, x)
	}
	return res
}

func main() {
	http.HandleFunc(
		"/factorial", func(writer http.ResponseWriter, request *http.Request) {
			query := request.URL.Query()
			numberString, ok := query["number"]
			if !ok {
				http.Error(writer, "number param is required", http.StatusBadRequest)
				return
			}
			number, err := strconv.Atoi(numberString[0])
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}

			start := time.Now()
			res := factorial(big.NewInt(int64(number)))
			duration := time.Now().Sub(start)

			bytes := res.Bytes()

			fmt.Fprintf(
				writer,
				"Duration = %s\nSize = %d\nResult = %s",
				duration.String(),
				len(bytes),
				res.String(),
			)
		},
	)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
