#!/usr/bin/env ruby

def fn(dpath)
  puts dpath
  if File.exist?(dpath+'/.grive')
    puts "---"
    Dir.chdir(dpath) do
      puts `drive quota`
      `grive`
    end
    if File.exist?(dpath+'/.trash')
      puts "\n\n================\n"
      puts "=== WARNING! ===\n"
      puts "= trash exist. =\n"
      puts "================\n"
    end
  else
    newpath = File.absolute_path(dpath+'/..')
    if newpath != '/'
      fn(newpath)
    end
  end
end

if ARGV[0] = '-e'
  `drive emptytrash`
end

fn `pwd`.chomp
