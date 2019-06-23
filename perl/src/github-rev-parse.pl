#!perl

use strict;
use warnings;
use utf8;
use Getopt::Long qw(:config posix_default no_ignore_case gnu_compat);
use HTTP::Tiny;
use JSON::PP qw/decode_json/;

sub usage {
    my ($msg) = @_;

    die <<"EOS";
ERROR: $msg
[usage]
  \$ perl github-rev-parse <org> <repo> <key (commit hash, branch, tag)>
  options:
    --token=github-token : pass the token of GitHub
EOS
}

sub main {
    my ($org, $repo, $key, $github_token) = @_;

    my %http_opt = (
        timeout => 5,
    );
    if ($github_token) {
        $http_opt{default_headers} = {
            Authorization => "token $github_token",
        };
    }

    my $res = HTTP::Tiny->new(%http_opt)->get("https://api.github.com/repos/${org}/${repo}/commits/${key}");
    if ($res->{success}) {
        my $body = decode_json($res->{content});
        print "$body->{sha}\n";
        return;
    }

    # there is no result that is matched
    exit(1);
}

my $github_token = '';
GetOptions(
  'token=s' => \$github_token,
);

my $org = shift or usage("the org name is missing");
my $repo = shift or usage("the repository name is missing");
my $key = shift or usage("the key is missing");

main($org, $repo, $key, $github_token);

__END__

