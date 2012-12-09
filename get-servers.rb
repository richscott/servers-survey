#! /usr/local/bin/ruby

require 'net/http'
require './thread-pool.rb'

Net::HTTP::read_timeout = 2.0

@headerRxs = [Regexp.new('server', Regexp::IGNORECASE),
              Regexp.new('x-powered-by', Regexp::IGNORECASE)]

# @headers should look like:
#  {'server': { 'apache': 100, 'IIS6': 80 },
#   'x-powered-by': {'foo': 80, 'bar': 43}
#  }
@headers = {}

def survey_site(website, name, count)
  printf "%5d. %-30s  %s\n", count, name, website
  uri = URI(website)
  begin
    res = Net::HTTP.get_response(uri)
  rescue Exception => e
    printf "Timeout Error (%s): %s\n", e.to_s, website
    return
  end

  if [200,301,302].none?{|code| code == res.code.to_i}
    puts "ERROR: couldn't query #{website} (#{name})"
    return
  end

  res.header.each_header {|hkey, hval|
    saveHeader = false
    if @headerRxs.none? {|re| re =~ hkey}  # ignore most headers
      next
    end

    #puts "I see header #{hkey}"
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
  if line =~ /^b/i
    break
  end
  next if line =~ /^\s*#/
  ary = line.split("\t")
  name= ary[0]
  website = ary[6]
  sites[website] = name
}
csvFile.close

$stdout.sync = true
$stderr.sync = true

pool = Pool.new(20)
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
