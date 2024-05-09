import string
import urllib.request


def analyze_text(url):
    with urllib.request.urlopen(url) as response:
        text = response.read().decode('utf-8').lower()
        hashmap = {char: {letter: 0 for letter in string.ascii_lowercase} for char in string.ascii_lowercase}

        for i, char in enumerate(text):
            if char.isalpha():
                next_index = i + 1
                while next_index < len(text) and not text[next_index].isalpha():
                    next_index += 1

                if next_index < len(text):
                    next_char = text[next_index]
                    hashmap[char][next_char] += 1

        return hashmap


def display(hashmap):
    user_input = input("Enter a letter: ").strip().lower()

    if len(user_input) == 1 and user_input.isalpha():
        output = hashmap.get(user_input)
        for key, value in output.items():
            print(f'{key}: {value}')
    else:
        print('Invalid value!')

# display(analyze_text("http://www.textfiles.com/internet/alt-news.txt"))
# display(analyze_text("http://www.textfiles.com/internet/dummy20.txt"))
# display(analyze_text("http://www.textfiles.com/programming/crenshawtut.txt"))