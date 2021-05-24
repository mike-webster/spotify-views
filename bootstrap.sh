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
    else 
        echo "incorrect go version, have: $installed; need: $version"; 
        echo "use `gvm list` to see if you have $version";
        echo "use `gvm install ${version#$rmv}` to install correct version";
        exit;
    fi
fi