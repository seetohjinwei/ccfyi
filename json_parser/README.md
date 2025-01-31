# JSON Parser / JQ

- https://codingchallenges.fyi/challenges/challenge-json-parser
- https://codingchallenges.fyi/challenges/challenge-jq/

```sh
# open the program in $(EDITOR)
make

# run the program on an input file
./file_parser.py test_json/full_suite/pass1.json

curl -s 'https://dummyjson.com/quotes?limit=2' | ./jq.py .
curl -s 'https://dummyjson.com/quotes?limit=2' | ./jq.py '.quotes'
curl -s 'https://dummyjson.com/quotes?limit=2' | ./jq.py '.quotes[].id'

curl -sL 'https://api.github.com/repos/CodingChallegesFYI/SharedSolutions/commits?per_page=3' | ./jq.py '.[0]'
curl -sL 'https://api.github.com/repos/CodingChallegesFYI/SharedSolutions/commits?per_page=3' | ccjq '.[0] | .commit.message'

# run all tests
make test
# run specific test function
make test test=tests.test_states.TestStringState.test_basic_string
# run entire test class
make test test=tests.test_states.TestStringState
make test test=tests.test_json_struct.TestJSONStruct_pretty_print
# run entire test module
make test test=tests.test_states
```
