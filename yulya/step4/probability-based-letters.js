const fs = require('fs');
const readline = require('readline');

if (process.argv.length < 4) {
  console.error('Example: node statistics-counter.js filename.txt letter');
  process.exit(1);
}

const filePath = process.argv[2];
const letter = process.argv[3];

if (letter.length !== 1 || !letter.match(/[a-z]/i)) {
  console.error('The letter argument should be a single alphabet character.');
  process.exit(1);
}

const dictionary = new Map();

function createDictionary (line, dictionary) {
  line = line.toLowerCase().normalize('NFKC').replace(/\n/g, ' ').replace(/\r/g, '').replace(/\s+/g, ' ').replace(/[^a-z ]/g, '');
  for (let i = 0; i < line.length - 1; i++) {
    let currentChar = line[i];
    let nextChar = line[i + 1];

    if (!dictionary.has(currentChar)) {
      dictionary.set(currentChar, new Map([[nextChar, 1]]));
    } else {
      let nextCharMap = dictionary.get(currentChar);
      if(!nextCharMap.has(nextChar)) {
        nextCharMap.set(nextChar, 1);
      } else {
        nextCharMap.set(nextChar, nextCharMap.get(nextChar) + 1);
      }
    }
  }
  return dictionary;
}

function weightedRandom (currentLetterMap) {
  const entries = Array.from(currentLetterMap.entries());

  let totalWeight = entries.reduce((total, entry) => total + entry[1], 0)
  let random = Math.random() * totalWeight;
  let cumulativeWeight = 0;

  for (let [letter, weight] of entries) {
    cumulativeWeight += weight;
    if (random < cumulativeWeight) {
      return letter;
    }
  }
}

function findMostProbableLetters (letter, dictionary, iterations = 200) {
  let currentLetter = letter;
  let mostProbableLetters = [];

  for(let i = 0; i < iterations; i++) {
    if (dictionary.has(currentLetter)) {
      const currentLetterMap = dictionary.get(currentLetter);
      let nextLetter = weightedRandom(currentLetterMap);

      mostProbableLetters.push(nextLetter)
      currentLetter = nextLetter;

    } else {
      console.log(`No entries found for '${currentLetter}'.`);
      break;
    }
  }
  const letterList = mostProbableLetters.join('')
  console.log(letterList);
}

try {
  fs.accessSync(filePath, fs.constants.F_OK | fs.constants.R_OK);
  const line = fs.readFileSync(filePath, 'utf-8');

  createDictionary(line, dictionary);

  findMostProbableLetters(letter, dictionary);
} catch (err) {
  console.error('Error handling the file', err);
}

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

const question1 = () => {
  return new Promise((resolve) => {
    rl.question(`Type next letter or 'quit' to exit : `, (answer) => {
      if (answer.toLowerCase() === 'quit') {
        console.log('Exiting application.');
        rl.close();
      } else {
        findMostProbableLetters(answer, dictionary);
        resolve(question1());
      }
    })
  })
};

const main = async () => {
  await question1();
};

main();
