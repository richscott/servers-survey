#! /usr/local/bin/ruby

require 'net/http'
require 'threadpool'

csvFile = File.open('fortune1000_companies.csv', 'r')

scanHeaders = ['server', 'x-powered-by']
#
# headers should look like:
#  {'server': { 'apache': 100, 'IIS6': 80 },
#   'x-powered-by': {'foo': 80, 'bar': 43}
#  }
headers = {}
count  = 0

$stdout.sync = true
$stderr.sync = true

sites = {}    # key is a URL, value is company name
csvFile.each { |line|
  next if line =~ /^\s*#/
  ary = line.split("\t")
  name= ary[0]
  website = ary[6]
  sites[website] = name
}

sites.each_pair{|website,name|
  uri = URI(website)
  begin
    res = Net::HTTP.get_response(uri)
  rescue Exception => e
    printf "Timeout Error (%s): %s\n", e.to_s, website
    next
  end

  if [200,301,302].include?(res.code.to_i)
    printf "%-30s  %s\n", name, website
    res.header.each_header {|hkey, hval|
      next if not scanHeaders.include?(hkey)  # ignore most headers

      if not headers.has_key?(hkey)
        headers[hkey] = {}
      end
      if headers[hkey].has_key?(hval)
        headers[hkey][hval] += 1
      else
        headers[hkey][hval] = 1
      end
    }
  else
    puts name + "   " + website  + '  ERROR: couldnt query site'
  end
  count += 1
  if count == 500
    break
  end
}
csvFile.close
puts ""

headers.keys.sort.each { |header|
  puts "HEADER: " + header
  headerValues = headers[header]
  headerValues.keys.each{|headerVal|
    printf "%5d  %s\n", headerValues[headerVal].to_s, headerVal
  }
}
