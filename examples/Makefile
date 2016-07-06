SAMPLES = $(wildcard *-sample.md)
EXE = $(subst -sample.md, , $(SAMPLES))
GO = caps.go comments.go deemph.go deflists.go
PDF = $(subst .md,.pdf, $(SAMPLES))

.PHONY: ${EXE}

all: ${EXE} ${SAMPLES} ${PDF} show

show:
	evince ${PDF}

${EXE}: ${GO}
	go build $@.go

%.pdf: %.md
	pandoc --filter ./$(subst -sample.md,,$<) $< -o $@ || touch $@

clean:
	$(RM) ${PDF}
