unit sub MAIN($start = 1, $stop = Inf);

my $active-lock = Lock.new;
my $primes-lock = Lock.new;

my @small-primes = (2,3,*+2... * > 997).grep: *.is-prime;

my @primes;
my @active;

($start...^$stop).race(degree => 3).map: -> $n {
  $active-lock.protect: { @active.push($n) };
  say "Testing : $n; Active: {@active.sort.gist}; Primes at {@primes.gist}";
  if ([1..$n].join ~ [$n^...1].join).&miller-rabin(1) {
	$primes-lock.protect: { @primes.push($n) };
  }
  $active-lock.protect: { @active.splice(@active.first($n, :k), 1) };
  say "Finished: $n; Active: {@active.sort.gist}; Primes at {@primes.gist}";
}

sub miller-rabin($n, $k) {
  return False if so $n %% any @small-primes;
  my ($s, $d) = find-s-and-d($n);
  for 1..$k {
	#my $a = (2..($n - 2)).pick;
	my $a = @small-primes.pick;
	my $x = expmod($a, $d, $n);
	my $y = 0;
	for 1..$s {
	  $y = ($x ** 2) mod $n.Int;
	  return False if $y == 1 and $x != 1 and $x != ($n - 1);
	  $x = $y;
 	}
	return False if $y != 1;
  }
  return True;
}

sub find-s-and-d($n) {
  my $d = $n - 1;
  my $s = 0;
  while ($d %% 2) {
	$d = $d div 2;
	$s++;
  }
  return $s, $d;
}
