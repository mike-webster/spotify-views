# TODO:

# Check dependencies
# - If something is missing, provide a link
# - If something is the wrong version, print a warning

echo 'installing spotify-views...'

echo 'checking for go:'
if ! command -v go  &> /dev/null
then 
    # go is missing
    echo "missing go installation, install gvm: https://github.com/moovweb/gvm#installing"
    exit
else 
    # check version
    rmv="go"
    version=$(cat .go-version)
    installed=$(go version)
    if [[ "$installed" == *"$version"* ]]; 
    then 
        echo "found $version"; 
        echo "go success!"
    else 
        echo "incorrect go version, have: $installed; need: $version"; 
        echo "use `gvm list` to see if you have $version";
        echo "use `gvm install ${version#$rmv}` to install correct version";
        exit;
    fi
fi

echo "checking for mysql..."
if ! command -v mysql  &> /dev/null
then 
    # missing
    echo 'missing mysql; please install mysql 5.7';
    exit;
else 
    # check version
    version=$(cat .mysql-version)
    installed=$(echo $(mysql --version) | sed  's/^.*Ver \([0-9]\.[0-9]\.[0-9].*\) for.*/\1/')

    if [[ "$installed" == *"$version"* ]]; 
    then 
        echo "found myql $version"; 
        echo "mysql success!"
    else 
        echo "incorrect go version, have: $installed; need: $version"; 
        echo "use `gvm list` to see if you have $version";
        echo "use `gvm install ${version#$rmv}` to install correct version";
        exit;
    fi
fi

echo 'building app...'
make dev