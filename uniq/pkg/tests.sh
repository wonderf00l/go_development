#!/bin/bash

test_file1=pkg/data-samples/input_1.txt
test_file2=pkg/data-samples/input_2.txt
test_file3=pkg/data-samples/unicode_1.txt
test_file4=pkg/data-samples/unicode_2.txt
test_file5=pkg/data-samples/sample.txt
test_file6=pkg/data-samples/empty.txt
test_file7=pkg/data-samples/newlines.txt

short_opts=("-d" "-u" "-c" "-i")  
files_list=($test_file1 $test_file2 $test_file3 $test_file4 $test_file5 $test_file6 $test_file7)

fieldStart=0
fieldEnd=4

charStart=0
charEnd=4

entrypoint=cmd/main.go

for (( offset=$fieldStart; offset<=$fieldEnd; offset++ ))
do
    for (( charOffset=$charStart; charOffset<=$charEnd; charOffset++ ))
        do
            for opt in ${short_opts[@]}; do 
                for test_file in ${files_list[@]}; do
                    diff <(go run $entrypoint $opt -f $offset -s $charOffset $test_file) <(uniq $opt -f $offset -s $charOffset $test_file) && echo "OK: $opt -f $offset -s $charOffset -- '$test_file' " || echo "FAIL: $opt -f $offset -s $charOffset -- '$test_file' "
            done
        done
    done
done

for (( offset=$fieldStart; offset<=$fieldEnd; offset+=2 ))
do
    for (( charOffset=$charStart; charOffset<=$charEnd; charOffset+=2 ))
        do
            for test_file in ${files_list[@]}; do
                diff <(go run $entrypoint -f $offset -s $charOffset $test_file) <(uniq -f $offset -s $charOffset $test_file) && echo "OK: -f $offset -s $charOffset -- '$test_file' " || echo "FAIL: -f $offset -s $charOffset -- '$test_file' "
        done
    done
done


for (( offset=$fieldStart; offset<=$fieldEnd; offset+=2 ))
do
    for (( charOffset=$charStart; charOffset<=$charEnd; charOffset+=2 ))
        do
            for test_file in ${files_list[@]}; do
                diff <(go run $entrypoint -d -c -f $offset -s $charOffset $test_file) <(uniq -d -c -f $offset -s $charOffset $test_file) && echo "OK: -d -c -f $offset -s $charOffset -- '$test_file' " || echo "FAIL: -d -c -f $offset -s $charOffset -- '$test_file' "
        done
    done
done

for (( offset=$fieldStart; offset<=$fieldEnd; offset+=2 ))
do
    for (( charOffset=$charStart; charOffset<=$charEnd; charOffset+=2 ))
        do
            for test_file in ${files_list[@]}; do
                diff <(go run $entrypoint -u -c -f $offset -s $charOffset $test_file) <(uniq -u -c -f $offset -s $charOffset $test_file) && echo "OK: -u -c -f $offset -s $charOffset -- '$test_file' " || echo "FAIL: -u -c -f $offset -s $charOffset -- '$test_file' "
        done
    done
done

for (( offset=$fieldStart; offset<=$fieldEnd; offset+=2 ))
do
    for (( charOffset=$charStart; charOffset<=$charEnd; charOffset+=2 ))
        do
            for test_file in ${files_list[@]}; do
                diff <(go run $entrypoint -u -i -f $offset -s $charOffset $test_file) <(uniq -u -i -f $offset -s $charOffset $test_file) && echo "OK: -u -i -f $offset -s $charOffset -- '$test_file' " || echo "FAIL:-u -i -f $offset -s $charOffset -- '$test_file' "
        done
    done
done

for (( offset=$fieldStart; offset<=$fieldEnd; offset+=2 ))
do
    for (( charOffset=$charStart; charOffset<=$charEnd; charOffset+=2 ))
        do
            for test_file in ${files_list[@]}; do
                diff <(go run $entrypoint -i -c -d -f $offset -s $charOffset $test_file) <(uniq -i -c -d -f $offset -s $charOffset $test_file) && echo "OK: -i -c -d -f $offset -s $charOffset -- '$test_file' " || echo "FAIL: -i -c -d -f $offset -s $charOffset -- '$test_file' "
        done
    done
done
