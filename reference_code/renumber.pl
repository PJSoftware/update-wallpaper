#!/usr/bin/perl

use strict;
use warnings;

my $rnCount = 0;

foreach my $jpg (<Spotlight/*.jpg>) {
  my ($folder,$file) = ($1,$2) if $jpg =~ m{(Spotlight)/(.+)};
  next if $file =~ m{^\d{8}-};

  my $tfile = $file;
  $tfile =~ s{[, &+;_]+}{-}g;
  $tfile =~ s{-[.]}{.}g;
  $tfile = lc($tfile);
  $tfile =~ s{-©-}{--}g;
  $tfile =~ s{---+}{--}g;
  $tfile =~ s{-+$}{}g;
  
  my $cmd = qq{git log --pretty=format:"\%ad" --date=short -- "$jpg"};
  my $date = `$cmd`;
  $date =~ s{-}{}g;
  $tfile = "$folder/$date-$tfile";

  print "Renaming to $tfile\n";
  my $git = qq{git mv "$jpg" "$tfile"};
  `$git`;
  $rnCount++;
}

if ($rnCount) {
  print "$rnCount files renamed; committing!";
  my $commitPush = qq{git commit -m "renamed by script" && git push};
  `$commitPush`;
}