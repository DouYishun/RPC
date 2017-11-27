#! /bin/bash
i=1
while(( $i <= 10 ))
do
    (go run client.go) &
    let i++
done
wait
