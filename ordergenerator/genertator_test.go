package ordergenerator

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	start := time.Now()
	orderGenerator := New()

	pdfBytes, err := orderGenerator.Generate(prepareFakeOrder(t))
	require.NoError(t, err)

	pdf, err := os.CreateTemp("", "test.pdf")
	require.NoError(t, err)

	_, err = pdf.Write(pdfBytes.Bytes())
	require.NoError(t, err)
	require.NoError(t, pdf.Close())

	duration := time.Since(start)
	t.Logf("Report took: %s\n", duration)

	openPdf(t, pdf.Name())
}

func prepareFakeOrder(t *testing.T) OrderArgs {
	t.Helper()

	orderArgs := OrderArgs{}
	err := gofakeit.Struct(&orderArgs)
	require.NoError(t, err)

	return orderArgs
}

func openPdf(t *testing.T, pdfPath string) {
	t.Helper()

	cmd := exec.Command("xdg-open", pdfPath)
	err := cmd.Run()
	require.NoError(t, err)
}
