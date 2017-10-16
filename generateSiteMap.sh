#! /bin/sh

# static
echo "https://debatehub.org/"
echo "https://debatehub.org/mission"
echo "https://debatehub.org/blog/contact"

# dynamic
for i in $(curl -s https://debatehub.org/debate_pages  | pup "tbody tr td p a" | grep href  | cut -d\" -f2 )
do
	echo "https://debatehub.org"$i
done



for i in $(curl -s https://debatehub.org/trends | pup 'div.media-body h4 a' | grep href | cut -d\" -f2)
do
	echo "https://debatehub.org"$i
done


