#/bin/sh
for d in */ ; do
            printf "run $d\r"
            rm -f $d/result.txt.gz

            TIME=$({ time whattle run -i $d/data.txt.gz -o $d/result.txt.gz -- $d/params.json; })
            DIFF=$( diff $d/result.txt.gz $d/expected.txt.gz )
            if [ "$DIFF" != "" ] ; then
                echo "$d: RESULTS DIFFERENT OH NO"
            else
                echo "$d: ${TIME:5}"
            
    fi
done

