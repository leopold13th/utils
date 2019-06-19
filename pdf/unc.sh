echo "Uncompressing..."
pdftk "$1" output uncompressed.pdf uncompress
echo "sed in the work..."
sed -f removewatermarks uncompressed.pdf > unwatermarked.pdf
echo "Compressing..."
pdftk unwatermarked.pdf output fixed.pdf compress
echo "All Done"
