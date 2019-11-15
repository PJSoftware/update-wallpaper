#!/perl -w

use strict;
use warnings;

use Cwd;
use Image::ExifTool;
use Digest::SHA qw{ sha256_hex };
use File::Copy;

our $VERSION = 0.1;

my $LOCAL = $ENV{LocalAppData};
my $ASSET_FOLDER = "$LOCAL/Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy/LocalState/Assets";
my $WALLPAPERS = 'C:/_WorkingSVN_/wallpaper/_current';

my %assets = ();
my %asset_size = ();

use constant FILESIZE => 7;

%Image::ExifTool::UserDefined::petesplace = (
    GROUPS => { 0 => 'XMP', 1 => 'XMP-petesplace', 2 => 'Image' },
    NAMESPACE => { 'petesplace' => 'http://petesplace.id.au' },
    WRITABLE => 'string',
    BadTimeFlag => {},
    );
%Image::ExifTool::UserDefined = (
      # new XMP namespaces (ie. XMP-xxx) must be added to the Main XMP table:
    'Image::ExifTool::XMP::Main' => {
        defra => {
            SubDirectory => {
                TagTable => 'Image::ExifTool::UserDefined::petesplace'
                },
            },
        }
    );

browse_assets();
identify_wp_sized();

compare_existing();
copy_new_images();

exit;

#################################################################

sub browse_assets {
    my $cwd = cwd();
    chdir $ASSET_FOLDER;
    foreach my $asset (glob q{*}) {
        $assets{$asset} = 1;
    }
    chdir $cwd;
    return;
}

sub identify_wp_sized {
    foreach my $file (sort keys %assets) {
        my $exif = Image::ExifTool::new();
        my $tags = $exif->ImageInfo("$ASSET_FOLDER/$file");
        my $size = $tags->{ImageSize};
        if (!defined $size || $size ne '1920x1080') {
            delete $assets{$file};
        }
    }
    return;
}

sub compare_existing {
    my $cwd = cwd();
    chdir $WALLPAPERS;

    my %by_size = ();
    my $filesize = 7;
    foreach my $jpg (glob q{*}) {
        my $size = (stat $jpg)[FILESIZE];
        $by_size{$size}{files}{$jpg} = 1;
        $by_size{$size}{flag} = 1;
    }

    foreach my $asset (sort keys %assets) {
        my $size = (stat "$ASSET_FOLDER/$asset")[FILESIZE];
        if (defined $by_size{$size}{flag}) {
            foreach my $jpg (sort keys %{$by_size{$size}{files}}) {
                if (compare_images($jpg,"$ASSET_FOLDER/$asset")) {
                    delete $assets{$asset};
                }
            }
        }
    }
    return;
}

sub compare_images {
    my ($file1,$file2) = @_;
    open my $FH1,'<',$file1  || return 0;
    open my $FH2,'<',$file2  || return 0;
    binmode $FH1 ;
    binmode $FH2 ;
    my $rv = 1;
    while ($rv && read $FH1,my $block1,8192) {
        read$FH2,my $block2,8192;
        $rv = (sha256_hex($block1) eq sha256_hex($block2));
        }
    close $FH1;
    close $FH2;
    return $rv;
    }

sub copy_new_images {
    foreach my $asset (sort keys %assets) {
        my $src = "$ASSET_FOLDER/$asset";
        my $trg = "$WALLPAPERS/$asset.jpg";
        print "Copying to $trg\n";
        print "OOPS\n" unless -f $src;
        copy($src, $trg);
        print "NOT COPIED\n" unless -f $trg;
    }
    return;
}
