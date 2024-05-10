import collections
import argparse

parser = argparse.ArgumentParser()
parser.add_argument('letter')
parser.add_argument('path')
args = parser.parse_args()

input_letter = args.letter
input_path = args.path

"""
Modify previous program in the way that it will generate a text (up to 200 characters long) in such way that next character will be the character with the most count of the current one, starting with input letter. If there are several letters with same count, it can choose any of them. If there is no such next letter, program finishes.

For example, given text "Hello, people!" and initial letter "e" it can generate "elelelel...."
"""


def count_letters_after(filename, letter):
    letter_count = collections.Counter()
    letter = letter.lower()

    with open(filename, 'r') as file:
        for line in file:
            for i, c in enumerate(line[:-1]):
                if c.lower() == letter:
                    if line[i + 1].isalpha():
                        next_letter = line[i + 1].lower()
                        letter_count[next_letter] += 1
    return letter_count


def generate_text(filename, initial_letter, max_length=200):
    # initial text with letter in args
    text = initial_letter.lower()

    # lower() input letter from args
    current_letter = initial_letter.lower()

    # count leter after, using previous code
    counts_after_letter = count_letters_after(filename, initial_letter)

    while len(text) < max_length:
        next_letters = [letter for letter, count in counts_after_letter.items() if letter != current_letter]
        if not next_letters:
            break
        next_letter = max(next_letters, key=counts_after_letter.get)
        print(next_letter)
        text += next_letter
        current_letter = next_letter

    return text


generated_text = generate_text(input_path, input_letter, max_length=200)
print(generated_text)
