SAMPLES = $(wildcard *-sample.md)
EXE = $(subst -sample.md, , $(SAMPLES))
# GO = caps.go comments.go deemph.go deflists.go
GO = $(subst -sample.md,.go, $(SAMPLES))
PDF = $(subst .md,.pdf, $(SAMPLES))

.PHONY: ${EXE}

all: ${EXE} ${SAMPLES} ${PDF}

show:
	evince ${PDF}

${EXE}: ${GO}
	go build $@.go

%.pdf: %.md
	pandoc --filter ./$(subst -sample.md,,$<) $< -o $@
	
clean:
	$(RM) ${PDF}
