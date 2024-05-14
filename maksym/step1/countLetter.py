import collections

filename = 'file.txt'


def count_letters(filename):
    letter_count = collections.Counter()

    with open(filename, 'r') as file:
        for line in file:
            letter_count.update(c.lower() for c in line if c.isalpha())
    return sorted(letter_count.items())


counts = count_letters(filename)
for letter, count in counts:
    print(f'{letter}: {count}')