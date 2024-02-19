package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/night1010/everhealth/entity"
	"github.com/TigorLazuardi/tanggal"
	"github.com/jung-kurt/gofpdf"
)

func CreatePrescriptionLeave(Prescription string, user *entity.Profile, doctor *entity.DoctorProfile) (io.Reader, error) {
	var result bytes.Buffer
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 15)
	pdf.Image("./img/everhealth.png", 11, 10, 80, 0, false, "", 0, "")
	pdf.Cell(0, 30, "-----------------------------------------------------------------------------------------------------------")
	pdf.Ln(20)

	publicationDate, err := tanggal.Papar(time.Now(), "", tanggal.WIB)
	if err != nil {
		return nil, err
	}

	format := []tanggal.Format{
		tanggal.Hari,
		tanggal.NamaBulan,
		tanggal.Tahun,
	}
	pdf.SetFont("Arial", "", 16)
	pdf.Write(0, fmt.Sprintf("Tanggal Resep Obat dibuat: %v", publicationDate.Format(" ", format))) // Add underlined text
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 16)
	pdf.Cell(0, 0, fmt.Sprintf("Konsultasi dengan Dr. %v", doctor.Profile.Name))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 16)
	pdf.Cell(0, 0, fmt.Sprintf("Nama Pasien: %v", user.Name))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Resep Obat: ")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 15)
	pdf.Cell(0, 10, Prescription)
	pdf.Ln(10)

	err = pdf.Output(&result)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(&result)
	return reader, nil
}
