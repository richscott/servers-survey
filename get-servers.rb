#! /usr/local/bin/ruby

require 'open-uri'
require './thread-pool.rb' # from http://burgestrand.se/code/ruby-thread-pool/

@headerRxs = [Regexp.new('server', Regexp::IGNORECASE),
              Regexp.new('x-powered-by', Regexp::IGNORECASE)]

# @headers should look like:
#  {'server': { 'apache': 100, 'IIS6': 80 },
#   'x-powered-by': {'foo': 80, 'bar': 43}
#  }
@headers = {}

def survey_site(website, name, count)
  printf "%5d. %-30s  %s\n", count, name, website
  begin
    http = open(website, 'r', :read_timeout => 5.0, :redirect => true)
  rescue Exception => e
    printf "Timeout Error for %s (%s)\n", website, e.to_s
    return
  end

  if [200,301,302].none?{|code| code == http.status[0].to_i}
    puts "ERROR: couldn't query #{website} (#{name}) - status was " + http.status.join(' ')
    return
  end

  http.meta.each {|hkey, hval|
    if @headerRxs.none? {|re| re =~ hkey}  # ignore most headers
      next
    end

    if not @headers.has_key?(hkey)
      @headers[hkey] = {}
    end
    if @headers[hkey].has_key?(hval)
      @headers[hkey][hval] += 1
    else
      @headers[hkey][hval] = 1
    end
  }
end

sites = {}    # key is a URL, value is company name
csvFile = File.open('fortune1000_companies.csv', 'r')
csvFile.each { |line|
  next if line =~ /^\s*#/
  ary = line.split("\t")
  name= ary[0]
  website = ary[6]
  sites[website] = name
}
csvFile.close

$stdout.sync = true
$stderr.sync = true

pool = Pool.new(10)
@count  = 1
sites.each_pair{|website,name|
  pool.schedule {
    survey_site(website, name, @count)
    @count += 1
  }
}

pool.shutdown

@headers.keys.sort.each { |header|
  puts "\nHEADER: " + header
  headerValues = @headers[header]
  headerValues.keys.sort.each{|headerVal|
    printf "%5d  %s\n", headerValues[headerVal].to_s, headerVal
  }
}
#at_exit { pool.shutdown }
