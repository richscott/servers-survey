#! /usr/local/bin/ruby

require 'net/http'
require 'threadpool'


@scanHeaders = ['server', 'x-powered-by']
#
# @headers should look like:
#  {'server': { 'apache': 100, 'IIS6': 80 },
#   'x-powered-by': {'foo': 80, 'bar': 43}
#  }
@headers = {}

$stdout.sync = true
$stderr.sync = true

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

def survey_site(website, name)
  uri = URI(website)
  begin
    res = Net::HTTP.get_response(uri)
  rescue Exception => e
    printf "Timeout Error (%s): %s\n", e.to_s, website
    return
  end

  if [200,301,302].include?(res.code.to_i)
    printf "%-30s  %s\n", name, website
    res.header.each_header {|hkey, hval|
      if not @scanHeaders.include?(hkey)  # ignore most headers
        return
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
  else
    puts name + "   " + website  + '  ERROR: couldnt query site'
  end
end

count  = 0
sites.each_pair{|website,name|
  survey_site(website, name)
  count += 1
  if count == 20
    break
  end
}
puts ""

@headers.keys.sort.each { |header|
  puts "HEADER: " + header
  headerValues = @headers[header]
  headerValues.keys.each{|headerVal|
    printf "%5d  %s\n", headerValues[headerVal].to_s, headerVal
  }
}
