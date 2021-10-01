package main

import (
	"fmt"
	"time"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	start := time.Now()

	db := core.Connect()

	db.Store("mytext", `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris id eros et risus mollis fringilla at non diam. Etiam nec malesuada felis. Maecenas in leo vitae augue commodo semper nec eu diam. Donec at justo velit. Ut sagittis enim sed neque vehicula scelerisque. Vestibulum sollicitudin tellus sit amet mauris vulputate, ac tempus massa feugiat. Aliquam lectus ante, aliquet vel facilisis a, mattis id mi. Suspendisse eget condimentum dolor, sit amet mattis nunc. Quisque dapibus ipsum quis ante vehicula venenatis. Mauris ac velit non augue efficitur consectetur non sit amet urna. Sed facilisis tortor vitae posuere congue. Etiam eu laoreet turpis. Sed non quam ac felis rutrum dictum. Etiam a porttitor tellus.
	
	In in enim urna. Curabitur luctus turpis in ornare maximus. Maecenas sollicitudin, ipsum ut semper placerat, dolor eros congue ante, vel dapibus mauris leo vel velit. Sed imperdiet maximus erat sed volutpat. Mauris et aliquam lorem. Etiam dictum non leo vel commodo. Ut elementum metus eu cursus cursus. Ut arcu odio, iaculis vitae porttitor sit amet, rutrum vel eros. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Praesent rutrum ut dui sed lobortis. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Pellentesque feugiat felis et fringilla porttitor.
	
	Nunc posuere pellentesque eros, quis venenatis nibh faucibus in. Praesent molestie diam et mauris cursus ullamcorper. Curabitur ex diam, lobortis vel mauris at, vestibulum tincidunt ex. Maecenas fermentum ipsum et mauris porttitor bibendum. Integer ultrices id enim placerat molestie. Nulla facilisi. Nulla eu tristique diam. Mauris eu erat eget augue aliquet dictum.
	
	Vestibulum sed dictum velit. Nulla nec tristique augue, sed pharetra magna. Donec nec magna et elit fringilla lobortis. Aliquam erat volutpat. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Duis euismod in purus a accumsan. Sed erat mauris, suscipit eu tempus vitae, lacinia quis ex. Maecenas ac massa eu mauris mattis vulputate. Donec eget dui vel leo porttitor fermentum quis eu purus. Mauris nisl mi, interdum vel turpis vel, sodales vehicula metus.
	
	Praesent dictum, massa ut rutrum aliquet, est lorem dignissim neque, volutpat sodales arcu elit eu arcu. Aenean vel egestas augue. Phasellus et pulvinar quam, a sagittis urna. Nunc molestie ipsum vitae orci egestas, non auctor nunc iaculis. Fusce eget leo augue. Nunc tincidunt, dui id eleifend ultrices, augue urna consequat ante, eget dignissim sem nibh vel purus. Praesent a tempor purus, id fringilla lorem. Cras consequat blandit sodales. Aenean feugiat accumsan lacus in semper. Curabitur vehicula turpis velit, vulputate viverra metus faucibus id. Pellentesque at luctus elit, nec eleifend nunc. Sed porta orci in nunc molestie consequat. Etiam in dui dolor. Morbi maximus tempor arcu ac malesuada. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Sed sed elit sapien.`)

	fmt.Println(db.Load("mytext"))
	db.Save()

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
