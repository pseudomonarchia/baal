#!/bin/sh

STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep ".go$")

if [[ "$STAGED_GO_FILES" = "" ]]; then
    exit 0
fi

PASS=true

go vet *.go
if [[ $? != 0 ]]; then
    PASS=false
fi

for FILE in $STAGED_GO_FILES
do
    if [[ $FILE == "vendor"* ]];then
        continue
    fi

    goimports -w $FILE
    if [[ $? != 0 ]]; then
        PASS=false
    fi

    golint "-set_exit_status" $FILE
    if [[ $? == 1 ]]; then
        PASS=false
    fi

    UNFORMATTED=$(gofmt -l $FILE)
    if [[ "$UNFORMATTED" != "" ]];then
        gofmt -w $PWD/$UNFORMATTED
        if [[ $? != 0 ]]; then
            PASS=false
        fi
    fi

    # git add $FILE
done

if ! $PASS; then
    printf "\033[31m FAILED \033[0m\n"
    exit 1
else
    printf "\033[32m SUCCEEDED \033[0m\n"
fi

exit 0
