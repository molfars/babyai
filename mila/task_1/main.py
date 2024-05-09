import urllib.request
import string


def analyze_text(url):
    with urllib.request.urlopen(url) as response:
        text = response.read().decode('utf-8').lower()

        hashmap = {char: 0 for char in string.ascii_lowercase}
        count = 0

        for char in text:
            if char.isalpha():
                hashmap[char] += 1
                count += 1

        for key, value in hashmap.items():
            print(f'{key}: {(value / count) * 100:.2f}%') if value else print(f'{key}: {value}')

        return hashmap

# print(analyze_text("http://www.textfiles.com/internet/alt-news.txt"))
# print(analyze_text("http://www.textfiles.com/internet/dummy20.txt"))
# print(analyze_text("http://www.textfiles.com/programming/crenshawtut.txt"))