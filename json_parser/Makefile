VENV=../venv
REQUIREMENTS=../requirements.txt
PYTHON=$(VENV)/bin/python3
PIP=$(VENV)/bin/pip
EDITOR=vim
VENV_ACTIVATE=. $(VENV)/bin/activate

.ONESHELL:

.PHONY: all
all: edit

.PHONY: edit
edit: $(VENV)/bin/activate
	$(VENV_ACTIVATE) && $(EDITOR)

.PHONY: test
test: $(VENV)/bin/activate
	if [ -z "$(test)" ]; then \
		$(VENV_ACTIVATE) && $(PYTHON) -m unittest discover tests; \
	else \
		$(VENV_ACTIVATE) && $(PYTHON) -m unittest $(test); \
	fi

.PHONY: env
env: $(VENV)/bin/activate

$(VENV)/bin/activate: $(REQUIREMENTS)
	rm -rf $(VENV)
	python3 -m venv $(VENV)
	$(VENV_ACTIVATE) && $(PIP) install -r $(REQUIREMENTS)

.PHONY: venv
venv: $(VENV)/bin/activate

.PHONY: clean
clean:
	rm -rf __pycache__
	rm -rf $(VENV)

.PHONY: lint
lint: $(VENV)/bin/activate
	ruff check
	ruff format
