import os
import random
import urllib.request, urllib.error

ALLOWED_CHARACTERS = 'abcdefghijklmnopqrstuvwxyz '


def analyze_text(url):
    try:
        urllib.request.urlopen(url)
    except urllib.error.URLError:
        raise ValueError("Provided URL is invalid.")

    if not os.path.splitext(url)[1] == ".txt":
        raise ValueError("Provided URL doesn't point to a .txt file")

    with urllib.request.urlopen(url) as response:
        text = response.read().decode('utf-8').lower()

        global count
        count = 0

        hashmap = {char: {letter: 0 for letter in ALLOWED_CHARACTERS} for char in ALLOWED_CHARACTERS}

        for i, char in enumerate(text):
            if char in ALLOWED_CHARACTERS:
                next_index = i + 1
                count += 1

                while next_index < len(text) and not text[next_index] in ALLOWED_CHARACTERS:
                    next_index += 1

                if next_index < len(text):
                    next_char = text[next_index]
                    hashmap[char][next_char] += 1

        return hashmap


def display(hashmap):
    user_input = input("Enter a letter: ").strip().lower()

    if len(user_input) == 1 and user_input in ALLOWED_CHARACTERS:
        return generate_text(hashmap, user_input, 200)
    else:
        raise ValueError('Invalid value!')


def generate_text(hashmap, input, given_length):
    generated = input
    current_char = input

    for char in hashmap:
        # sum of the values for the current char's hashmap
        total_sum = sum(hashmap[char].values())

        # supposed sum of the values for the current char's new scaled hashmap
        scaled_sum = round(total_sum / count) * given_length

        for next_char in hashmap[char]:
            if total_sum and scaled_sum:
                hashmap[char][next_char] = round((hashmap[char][next_char] / total_sum) * scaled_sum)

    while len(generated) != given_length:
        next_probabilities = hashmap[current_char]
        print(hashmap[current_char])
        next_char = random.choices(list(next_probabilities.keys()), list(next_probabilities.values()))[0]

        hashmap[current_char][next_char] -= 1

        if hashmap[current_char][next_char] == 0:
            del hashmap[current_char][next_char]

        generated += next_char
        current_char = next_char

    return generated

# print(display(analyze_text("http://www.textfiles.com/internet/alt-news.txt")))
# print(display(analyze_text("http://www.textfiles.com/internet/dummy20.txt")))
# print(display(analyze_text("http://www.textfiles.com/programming/crenshawtut.txt")))
