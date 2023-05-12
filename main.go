package main

import (
	"awesomeProject1/model"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	batchSize := 500

	filename := "meta.tsv"
	fileReader, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERR: ", err)
	}
	scanner := bufio.NewScanner(fileReader)
	var head = model.Header{}
	var serviceMeta = map[int]model.ServiceMeta{}
	var songs = map[int]model.SongUsage{}
	var DSPReleaseID string
	var batchCounter = 0
	var outputSongUsages []model.SongUsage
	var rowData int

	//LOOP BARIS
	for scanner.Scan() {
		rowData++
		line := scanner.Text()
		var header string
		var arrayCount int
		var insertDb bool
		var headerConstructed bool

		var song = model.SongUsage{}
		var service = model.ServiceMeta{}

		// Mapping Data
		mappingData(line, rowData)
		//LOOP KOLOM
		for _, word := range strings.Split(line, "\t") {
			if word == "" {
				arrayCount++
				continue
			}

			// CHECK HEAD
			if header == "" {
				switch word {
				case "HEAD":
					header = "HEADER"
				case "SY02.03":
					header = "HEADER_USAGE"
				case "SR03.01":
					header = "RELEASE"
					songs = map[int]model.SongUsage{}
					DSPReleaseID = ""
				case "SRFO":
					header = "FOOT"

					arrayCount++
					continue
				}

				//CHECK HEADER DATA
				if header == "HEADER" {
					switch arrayCount {
					case 1:
						head.MessageVersion = word
					case 2:
						head.Profile = word
					case 3:
						head.ProfileVersion = word
					case 4:
						head.MessageId = word
					case 5:
						head.DateTimeCreated, _ = time.Parse(time.RFC3339, word)
					case 6:
						head.FileNumber, _ = strconv.Atoi(word)
					case 7:
						head.NumberOfFiles, _ = strconv.Atoi(word)
					case 8:
						head.StartDate, _ = time.Parse("2006-01-02", word)
					case 9:
						head.EndDate, _ = time.Parse("2006-01-02", word)
					case 10:
						head.SenderPartyId = word
					case 11:
						head.SenderName = word
					case 12:
						head.ServiceDescription = word
					case 14:
						head.RecipientName = word
					}
					arrayCount++
				}

				//CHECK SY02.03
				if header == "HEADER_USAGE" {
					switch arrayCount {
					case 1:
						parseInt, _ := strconv.ParseInt(word, 0, 64)
						service.SummaryRecordID = int(parseInt)
					case 5:
						service.UseType = word
					case 7:
						service.ServiceDescription = word
					case 8:
						usages, _ := strconv.ParseFloat(word, 64)
						service.Usages = usages
					case 11:
						netRevenue, _ := strconv.ParseFloat(word, 64)
						service.NetRevenue = netRevenue
					}
					if service.SummaryRecordID != 0 {
						serviceMeta[service.SummaryRecordID] = service
					}
					fmt.Println("SERVICE : ", service)

					arrayCount++
				}

				//CHECK RE01
				if header == "RELEASE" {
					switch arrayCount {
					case 3:
						DSPReleaseID = word
					}

					arrayCount++
				}

				//CHECK AS02.02 DATA
				if header == "ARTISTSONG" {
					switch arrayCount {
					case 1:
						song.BlockId, _ = strconv.Atoi(word)
					case 2:
						song.ResourceReference, _ = strconv.Atoi(word)
					case 3:
						song.DspResourceId = word
					case 4:
						song.ISRC = word
					case 5:
						song.Title = word
					case 7:
						song.DisplayArtistName = word
					case 8:
						song.Duration = word
					case 9:
						song.ResourceType = word
					case 12:
						song.ComposerAuthor = strings.Split(word, "|")
					}

					if song.ResourceReference != 0 {
						songs[song.ResourceReference] = song
					}
					arrayCount++
				}

				//CHECK SU02 DATA
				if header == "SALESUSAGE" {
					switch arrayCount {
					case 5:
						i, _ := strconv.Atoi(word)
						song = songs[i]
					case 6:
						song.IsRoyaltyBearing, _ = strconv.ParseBool(word)
					case 7:
						song.NumberOfStreams, _ = strconv.ParseFloat(word, 64)
						song.DSPReleaseID = DSPReleaseID
						insertDb = true
					}

					if song.ResourceReference != 0 {
						songs[song.ResourceReference] = song
					}
					arrayCount++
				}
				song.IndexNumber = fmt.Sprint(rowData)
			}

			if headerConstructed {
				fmt.Println("HEADER TRUE")
				//_, headerErr := headerConstructedCallback(head)
				//if headerErr != nil {
				//	return songProcessed, headerErr
				//}
			}

			if insertDb {
				batchCounter++
				csong := songs[song.ResourceReference]
				outputSongUsages = append(outputSongUsages, csong)
			}

			if batchCounter == batchSize {
				fmt.Println("BATCH => SongUsage")
				fmt.Println("OUTPUT : ", outputSongUsages)
				//err = songUsageConstructedCallback(head, outputSongUsages)
				//if err != nil {
				//	return songProcessed, err
				//}
				//songProcessed += len(outputSongUsages)
				batchCounter = 0
				outputSongUsages = []model.SongUsage{}
			}
		}

		//publish last batch
		if len(outputSongUsages) > 0 {
			fmt.Println("PUBLISH LAST BATCH")
			fmt.Println("OUTPUT : ", outputSongUsages)

			//err = songUsageConstructedCallback(head, outputSongUsages)
			//if err != nil {
			//	return songProcessed, err
			//}
			//songProcessed += len(outputSongUsages)
		}
	}

	//fmt.Println("SERVICE : ", serviceMeta)
}

func mappingData(line string, lineNumber int) {
	var serviceMeta []model.ServiceMeta
	row := strings.Split(line, "\t")
	for _, word := range strings.Split(line, "\t") {
		if word == "HEAD" {
			mappingHeader(row)
		} else if word == "SY02.03" {
			service := mappingService(row)
			serviceMeta = append(serviceMeta, service)
		} else if word == "SR03.01" {
			mappingSongUsages(row, lineNumber)
		} else if word == "SRFO" {

		}
	}
}

func mappingHeader(row []string) {
	dateTime, _ := time.Parse(time.RFC3339, row[5])
	fileNumber, _ := strconv.Atoi(row[6])
	numberOfFiles, _ := strconv.Atoi(row[7])
	startDate, _ := time.Parse("2006-01-02", row[8])
	endDate, _ := time.Parse("2006-01-02", row[9])
	var head = model.Header{
		HeaderID:           0,
		CatalogueHeaderID:  0,
		PublisherID:        0,
		MessageVersion:     row[1],
		Profile:            row[2],
		ProfileVersion:     row[3],
		MessageId:          row[4],
		DateTimeCreated:    dateTime,
		FileNumber:         fileNumber,
		NumberOfFiles:      numberOfFiles,
		StartDate:          startDate,
		EndDate:            endDate,
		SenderPartyId:      row[10],
		SenderName:         row[11],
		ServiceDescription: row[12],
		RecipientName:      row[14],
	}

	fmt.Println("HEAD : ", head)
}

func mappingService(row []string) (service model.ServiceMeta) {
	summaryRecordID, _ := strconv.ParseInt(row[1], 0, 64)
	usage, _ := strconv.ParseFloat(row[8], 64)
	netRevenue, _ := strconv.ParseFloat(row[1], 64)

	service = model.ServiceMeta{
		SummaryRecordID:    int(summaryRecordID),
		CommercialMode:     row[4],
		UseType:            row[5],
		ServiceDescription: row[7],
		Usages:             usage,
		NetRevenue:         netRevenue,
	}
	fmt.Printf("SERVICE %+v", service)

	return
}

func mappingSongUsages(row []string, lineNumber int) {
	numOfStream, _ := strconv.ParseFloat(row[32], 64)
	song := model.SongUsage{
		DspResourceId:     row[2],
		ISRC:              row[3],
		Title:             row[4],
		SubTitle:          row[5],
		DisplayArtistName: row[6],
		Duration:          row[8],
		ResourceType:      row[9],
		NumberOfStreams:   numOfStream,
		IndexNumber:       fmt.Sprint(lineNumber),
		Revenue:           0,
	}
	fmt.Printf("SONG :%+v", song)
}
