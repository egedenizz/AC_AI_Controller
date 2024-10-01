package speech

import (
	"bytes"

	"encoding/binary"
	"fmt"

	"os"

	"log"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

func ConvertToWav() []byte {

	rawData, err := os.ReadFile("output.raw")
	if err != nil {
		fmt.Println("Error reading raw file:", err)

	}

	data16 := make([]int16, len(rawData)/2)

	if err := binary.Read(bytes.NewReader(rawData), binary.LittleEndian, &data16); err != nil {
		fmt.Println("Error converting raw data:", err)

	}

	data := make([]int, len(data16))
	for i, v := range data16 {
		data[i] = int(v)
	}

	buffer := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  44100,
			NumChannels: 1,
		},
		Data: data,
	}

	outFile, err := os.Create("output.wav")
	if err != nil {
		fmt.Println("Error creating wav file:", err)

	}
	defer outFile.Close()

	encoder := wav.NewEncoder(outFile, buffer.Format.SampleRate, 16, buffer.Format.NumChannels, 1)
	if err := encoder.Write(buffer); err != nil {
		fmt.Println("Error writing wav file:", err)

	}

	if err := encoder.Close(); err != nil {
		fmt.Println("Error closing wav encoder:", err)
	}

	audioBytes, err := os.ReadFile("output.wav")
	if err != nil {
		fmt.Println("Error reading")
	}

	err = os.Remove("output.raw")
	if err != nil {

		fmt.Println("Error deleting file:", err)

	}
	err = os.Remove("output.wav")
	if err != nil {

		fmt.Println("Error deleting file:", err)

	}

	return audioBytes

}

func RecordVoice() {

	portaudio.Initialize()
	defer portaudio.Terminate()

	const sampleRate = 44100
	const seconds = 3
	buf := make([]int16, sampleRate*seconds)

	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(buf), buf)
	if err != nil {
		log.Fatalf("Failed to open PortAudio stream: %v", err)
	}
	defer stream.Close()

	fmt.Println("Recording...")
	if err := stream.Start(); err != nil {
		log.Fatalf("Failed to start PortAudio stream: %v", err)
	}
	if err := stream.Read(); err != nil {
		log.Fatalf("Failed to read from PortAudio stream: %v", err)
	}
	if err := stream.Stop(); err != nil {
		log.Fatalf("Failed to stop PortAudio stream: %v", err)
	}
	fmt.Println("Recording finished")

	file, err := os.Create("output.raw")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()
	binary.Write(file, binary.LittleEndian, buf)

}
