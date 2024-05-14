import collections
import argparse

parser = argparse.ArgumentParser()
parser.add_argument('letter')
parser.add_argument('path')
args = parser.parse_args()

input_letter = args.letter
input_path = args.path

"""
Count the number of letter after letter in parameter
Usage: python3 countProbability.py <letter>
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


counts_after_letter = count_letters_after(input_path, input_letter)

for letter, count in sorted(counts_after_letter.items()):
    print(f'{letter}: {count}')


