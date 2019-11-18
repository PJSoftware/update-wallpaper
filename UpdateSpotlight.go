package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const sourceFolder = "Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy/LocalState/Assets"
const targetFolder = "C:/_WorkingSVN_/wallpaper/1-screen/win10-spotlight"

var localAppData = os.Getenv("LOCALAPPDATA")
var sourcePath = localAppData + "/" + sourceFolder

var assetBySize map[int64]map[string]string

func main() {
	browseAssets(sourcePath)
	// scanExisting(targetFolder)
	// compare_existing()
	// copy_new_images()

	// summary()
	// pause()
}

func browseAssets(sourcePath string) {
	assetsFound := 0
	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	assetBySize = make(map[int64]map[string]string)
	for _, file := range files {
		filePath := sourcePath + "/" + file.Name()
		if isWallpaper(filePath, 1920, 1080) {
			fileSize := file.Size()
			if _, ok := assetBySize[fileSize]; !ok {
				assetBySize[fileSize] = make(map[string]string)
			}
			assetBySize[fileSize][filePath] = md5String(filepath)
			assetsFound++
		}
	}

	fmt.Printf("%d Spotlight images found\n", assetsFound)
}

func md5String(filePath string) string {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		hash := md5.New()
		if _, err := io.Copy(hash, file); err == nil {
			return hex.EncodeToString(hash.Sum(nil))
		}
	}

	return ""
}

func isWallpaper(filePath string, width, height int) bool {
	asset, err := os.Open(filePath)
	if err != nil {
		return false // Cannot read, so not interested in it
	}
	defer asset.Close()

	image, err := jpeg.DecodeConfig(asset)
	if err != nil {
		return false // Not a JPEG, so not interested in it
	}

	if image.Width != width || image.Height != height {
		return false
	}

	return true
}

func scanExisting(targetPath string) {
	wpFound := 0
	matchesFound := 0
	files, err := ioutil.ReadDir(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := targetPath + "/" + file.Name()
		fileSize := file.Size()
		if _, ok := assetBySize[fileSize]; ok {
			for assetName, assetHash := range assetBySize[fileSize] {
				wpHash := md5String(filePath)
			}
		}

	}

	fmt.Printf("%d Spotlight images found\n", assetsFound)

}

// sub browse_assets {
//     my $cwd = cwd();
//     chdir $ASSET_FOLDER;
//     foreach my $asset (glob q{*}) {
//         $assets{$asset} = 1;
//         $IN_ASSETS++;
//     }
//     chdir $cwd;
//     return;
// }

// #!/perl -w

// use strict;
// use warnings;

// use Cwd;
// use Image::ExifTool;
// use Digest::SHA qw{ sha256_hex };
// use File::Copy;
// use Readonly;

// our $VERSION = 0.1;

// my $LOCAL = $ENV{LocalAppData};
// my $ASSET_FOLDER = "$LOCAL/Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy/LocalState/Assets";
// my $WALLPAPERS = 'C:/_WorkingSVN_/wallpaper/1-screen/win10-spotlight';

// abort("Asset folder not found: $ASSET_FOLDER") unless -d $ASSET_FOLDER;
// abort("Wallpaper folder not found: $WALLPAPERS") unless -d $WALLPAPERS;

// my %assets = ();
// my %asset_size = ();

// Readonly my $FILESIZE_INDEX => 7;

// no warnings;    # Suppress "used only once" warning
// %Image::ExifTool::UserDefined::petesplace = (
//     GROUPS => { 0 => 'XMP', 1 => 'XMP-petesplace', 2 => 'Image' },
//     NAMESPACE => { 'petesplace' => 'http://petesplace.id.au' },
//     WRITABLE => 'string',
//     BadTimeFlag => {},
//     );
// %Image::ExifTool::UserDefined = (
//       # new XMP namespaces (ie. XMP-xxx) must be added to the Main XMP table:
//     'Image::ExifTool::XMP::Main' => {
//         defra => {
//             SubDirectory => {
//                 TagTable => 'Image::ExifTool::UserDefined::petesplace'
//                 },
//             },
//         }
//     );
// use warnings;

// my ($IN_ASSETS, $WP_ASSETS, $IN_WALLPAPER, $COPY_COUNT) = (0, 0, 0, 0);
// exit;

// #################################################################

// sub identify_wp_sized {
//     $WP_ASSETS = $IN_ASSETS;
//     foreach my $file (sort keys %assets) {
//         my $exif = Image::ExifTool::new();
//         my $tags = $exif->ImageInfo("$ASSET_FOLDER/$file");
//         my $size = $tags->{ImageSize};
//         if (!defined $size || $size ne '1920x1080') {
//             delete $assets{$file};
//             $WP_ASSETS--;
//         }
//     }
//     return;
// }

// sub compare_existing {
//     my $cwd = cwd();
//     chdir $WALLPAPERS;

//     my %by_size = ();
//     foreach my $jpg (glob q{*}) {
//         my $size = (stat $jpg)[$FILESIZE_INDEX];
//         $by_size{$size}{files}{$jpg} = 1;
//         $by_size{$size}{flag} = 1;
//         $IN_WALLPAPER++;
//     }

//     foreach my $asset (sort keys %assets) {
//         my $size = (stat "$ASSET_FOLDER/$asset")[$FILESIZE_INDEX];
//         if (defined $by_size{$size}{flag}) {
//             foreach my $jpg (sort keys %{$by_size{$size}{files}}) {
//                 if (compare_images($jpg,"$ASSET_FOLDER/$asset")) {
//                     delete $assets{$asset};
//                 }
//             }
//         }
//     }
//     return;
// }

// sub compare_images {
//     my ($file1,$file2) = @_;
//     open my $FH1,'<',$file1  || return 0;
//     open my $FH2,'<',$file2  || return 0;
//     binmode $FH1 ;
//     binmode $FH2 ;
//     my $rv = 1;
//     Readonly my $BLOCKSIZE => 8192;
//     my ($block1, $block2);
//     while ($rv && read $FH1, $block1, $BLOCKSIZE) {
//         read $FH2, $block2, $BLOCKSIZE;
//         $rv = (sha256_hex($block1) eq sha256_hex($block2));
//         }
//     close $FH1;
//     close $FH2;
//     return $rv;
//     }

// sub copy_new_images {
//     foreach my $asset (sort keys %assets) {
//         my $src = "$ASSET_FOLDER/$asset";
//         my $trg = "$WALLPAPERS/ZZZ_Unsorted_$asset.jpg";
//         print "Copying to $trg\n";
//         print "OOPS: $src does not exist\n" unless -f $src;
//         copy($src, $trg);
//         if (-f $trg) {
//             $COPY_COUNT++;
//         } else {
//             print "NOT COPIED: $trg\n" unless -f $trg;
//         }
//     }
//     return;
// }

// sub summary {
//     print "\n";
//     print "Identified $WP_ASSETS HD Wallpapers (of $IN_ASSETS files)\n";
//     print "Compared with $IN_WALLPAPER existing Spotlight Wallpapers\n";
//     print "Copied $COPY_COUNT files which did not already exist\n";
// }

// sub pause {
//     my ($msg) = @_;
//     $msg //= 'Press [Enter] to continue...';

//     print "$msg\n";
//     <STDIN>;
// }

// sub abort {
//     my ($msg) = @_;
//     pause($msg);
//     die $msg;
//     }
