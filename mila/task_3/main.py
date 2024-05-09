import random
import urllib.request

ALLOWED_CHARACTERS = 'abcdefghijklmnopqrstuvwxyz '


def analyze_text(url):
    with urllib.request.urlopen(url) as response:
        text = response.read().decode('utf-8').lower()
        hashmap = {char: {letter: 0 for letter in ALLOWED_CHARACTERS} for char in ALLOWED_CHARACTERS}

        for i, char in enumerate(text):
            if char in ALLOWED_CHARACTERS:
                next_index = i + 1
                while next_index < len(text) and not text[next_index] in ALLOWED_CHARACTERS:
                    next_index += 1

                if next_index < len(text):
                    next_char = text[next_index]
                    hashmap[char][next_char] += 1

        return hashmap


def display(hashmap):
    user_input = input("Enter a letter: ").strip().lower()

    if len(user_input) == 1 and user_input in ALLOWED_CHARACTERS:
        return generate_text(hashmap, user_input)
    else:
        return 'Invalid value!'


def generate_text(hashmap, user_input):
    generated = ' '
    next_char = user_input

    while len(generated) != 100:
        generated += next_char

        max_value = max(hashmap[next_char].values())
        keys = [key for key, value in hashmap[next_char].items() if value == max_value]
        # if several letters with maximum values => choose random letter
        next_char = random.choice(keys)

    return generated

# print(display(analyze_text("http://www.textfiles.com/internet/alt-news.txt")))
# print(display(analyze_text("http://www.textfiles.com/internet/dummy20.txt")))
# print(display(analyze_text("http://www.textfiles.com/programming/crenshawtut.txt")))