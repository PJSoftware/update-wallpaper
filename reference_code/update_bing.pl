#!/usr/bin/perl

use strict;
use warnings;

use File::Copy qw{ copy };

my $lad = $ENV{LOCALAPPDATA};
$lad =~ s{\\}{/}g;
my $bw = "$lad/Microsoft/BingWallpaperApp/WPImages";

my $copied = 0;

if (-d $bw) {
  my $trgt = "./Bing/$ENV{COMPUTERNAME}";
  mkdir "./Bing" unless -d "./Bing";
  mkdir $trgt unless -d $trgt;
  foreach my $f (<$bw/*.jpg>) {
    (my $fn = $f) =~ s{^$bw/}{};
    if (! -f "$trgt/$fn") {
      copy($f,"$trgt/$fn");
      print "Copied $fn\n";
      $copied++;
    }
  }
}

if ($copied) {
  my $commitPush = qq{git add . && git commit -m "add Bing wallpaper" && git push};
  `$commitPush`;
}
