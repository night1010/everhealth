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

func CreateSickLeave(startDate, endDate time.Time, user *entity.Profile, doctor *entity.DoctorProfile, diagnosa string) (io.Reader, error) {
	var result bytes.Buffer
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 15)
	pdf.Image("./img/everhealth.png", 11, 10, 80, 0, false, "", 0, "")
	pdf.Cell(0, 30, "-----------------------------------------------------------------------------------------------------------")
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Informasi Pasien")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 15)
	pdf.Cell(0, 10, fmt.Sprintf("Nama: %v", user.Name))
	pdf.Ln(10)

	format := []tanggal.Format{
		tanggal.Hari,
		tanggal.NamaBulan,
		tanggal.Tahun,
	}
	dob := user.Birthdate
	dobString, err := tanggal.Papar(dob, "", tanggal.WIB)
	if err != nil {
		return nil, err
	}
	pdf.SetFont("Arial", "", 15)
	pdf.Cell(0, 10, fmt.Sprintf("Tanggal Lahir: %v", dobString.Format(" ", format)))
	pdf.Ln(10)

	ageInt := age(dob, time.Now())
	pdf.SetFont("Arial", "", 15)
	pdf.Cell(0, 10, fmt.Sprintf("Umur: %v", ageInt))
	pdf.Ln(30)

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Catatan Dokter")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 16)
	publicationDate, err := tanggal.Papar(time.Now(), "", tanggal.WIB)
	if err != nil {
		return nil, err
	}
	pdf.Write(0, fmt.Sprintf("Tanggal Rekomendasi Terbit: %v", publicationDate.Format(" ", format))) // Add underlined text
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 16)
	pdf.Write(8, "Untuk Perhatian,")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 16)
	pdf.Write(8, fmt.Sprintf("Kami Informasikan bahwa pasien diatas kemungkinan didiagnosa: %v. Diagnosa tersebut diberikan atas gejala yang dijelaskan pasien pada saat konsultasi.", diagnosa))
	pdf.Ln(10)

	startDateIdn, err := tanggal.Papar(startDate, "", tanggal.WIB)
	if err != nil {
		return nil, err
	}
	endDateIdn, err := tanggal.Papar(endDate, "", tanggal.WIB)
	if err != nil {
		return nil, err
	}
	pdf.SetFont("Arial", "", 16)
	pdf.Write(8, fmt.Sprintf("Berdasarkan Keadaan tersebut, kami memberikan rekomendasi istirahat kepada pasien tersebut dari %v hingga %v", startDateIdn.Format(" ", format), endDateIdn.Format(" ", format)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 16)
	pdf.Write(8, "Demikian informasi dari kami, semoga dapat menjadi bahan pertimbangan")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 16)
	pdf.Cell(0, 0, "Konsultasi dengan")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 0, fmt.Sprintf("Dr. %v", doctor.Profile.Name))
	pdf.Ln(10)
	err = pdf.Output(&result)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(&result)
	return reader, nil
}

func age(birthdate, today time.Time) int {
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}
