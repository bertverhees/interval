package interval

/*
Reading the test-sets
The '========' stands for interval-range
Upper- en lower-Included is true
The '<' is for lowerUnbound, no lower-end
The '>' is for upperUnbound, no upper-end
The '*' is for Included is false (depending on the side of the string for upper en lower included)

The interval-string is divided in three parts by |,(or is empty)
First part indicates the lower-unbounded/included, the third part for the right side
It is allowed to have as many characters as is convenient on the first and the last part,
only the indicated characters are meaningful

It looks hard to understand and is maybe confusing, but read on,
So:
*|=================|  lowerIncluded is false, and lower  is 0
|=================|   lowerIncluded is true, and lower is 0
<|================| lowerUnbounded is true and lower is 0
<*|==============| lowerUnbounded is true, and lowerIncluded is false and lower is 0
*|---=================|  lowerIncluded is false, and lower  is 4
|---=================|   lowerIncluded is true, and lower is 4
|---<================| lowerUnbounded is true and lower is 4
<*|---==============| lowerUnbounded is true, and lowerIncluded is false and lower is 4

# On the upperside are mutatis mutandis the same rules

So, these two intervals below, although visible overlapping, do NOT overlap
<*|====--------------| 0,4
*|----=======|*       5,12
So it can be written like this, which is much more readable
So it can be written like this, which is much more readable
"<*|====--------------| 0,4"
"* |----=======|*"       5,12 or
" *|----=======|*"
Remember, this is only for creating tests-sets, it has nothing to do with a notation language of meaning outside these test-sets
*/

type testGeneral struct {
	i_interval_string string
	x_interval_string string
	counter           string
}

var testsGeneralSets = []testGeneral{
	{ //0
		i_interval_string: "  |=====---------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "0",
	},
	{ //1
		i_interval_string: "  |======--------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "1",
	},
	{ //2
		i_interval_string: "  |---=====------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "2",
	},
	{ //3
		i_interval_string: "  |-------=====--------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "3",
	},
	{ //4
		i_interval_string: "  |------========------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "4",
	},
	{ //5
		i_interval_string: "  |-----------=====----|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "5",
	},
	{ //6
		i_interval_string: "  |--------------======|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "6",
	},
	{ //7
		i_interval_string: "  |---------------=====|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "7",
	},
	{ //8
		i_interval_string: " *|======--------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "8",
	},
	{ //9
		i_interval_string: " *|------========------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "9",
	},
	{ //10
		i_interval_string: " *|--------------======|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "10",
	},
	{ //11
		i_interval_string: "  |======--------------|* ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "11",
	},
	{ //12
		i_interval_string: "  |------========------|* ",
		x_interval_string: "  |------========------|  ",
		counter:           "12",
	},
	{ //13
		i_interval_string: "  |--------------======|* ",
		x_interval_string: "  |------========------|  ",
		counter:           "13",
	},
	{ //14
		i_interval_string: " <|-------=====--------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "14",
	},
	{ //15
		i_interval_string: " <|------========------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "15",
	},
	{ //16
		i_interval_string: " <|-----------=====----|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "16",
	},
	{ //17
		i_interval_string: " <|--------------======|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "17",
	},
	{ //18
		i_interval_string: " <|---------------=====|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "18",
	},
	{ //19
		i_interval_string: "  |=====---------------|> ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "19",
	},
	{ //20
		i_interval_string: "  |======--------------|> ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "20",
	},
	{ //21
		i_interval_string: "  |---=====------------|> ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "21",
	},
	{ //22
		i_interval_string: "  |-------=====--------|> ",
		x_interval_string: "  |------========------|  ",
		counter:           "22",
	},
	{ //23
		i_interval_string: "  |------========------|> ",
		x_interval_string: "  |------========------|  ",
		counter:           "23",
	},
	{ //24
		i_interval_string: "  |======--------------|  ",
		x_interval_string: " *|------=======-------|  ",
		counter:           "24",
	},
	{ //25
		i_interval_string: "  |------========------|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "25",
	},
	{ //26
		i_interval_string: "  |--------------======|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "26",
	},
	{ //27
		i_interval_string: "  |======--------------|  ",
		x_interval_string: "  |------=======-------|*  ",
		counter:           "27",
	},
	{ //28
		i_interval_string: "  |------========------|  ",
		x_interval_string: "  |------========------|* ",
		counter:           "28",
	},
	{ //29
		i_interval_string: "  |--------------======|  ",
		x_interval_string: "  |------========------|*  ",
		counter:           "29",
	},
	{ //30
		i_interval_string: " *|======--------------|  ",
		x_interval_string: " *|------=======-------|  ",
		counter:           "30",
	},
	{ //31
		i_interval_string: " *|------========------|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "31",
	},
	{ //32
		i_interval_string: " *|--------------======|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "32",
	},
	{ //33
		i_interval_string: "  |======--------------|* ",
		x_interval_string: "  |------=======-------|* ",
		counter:           "33",
	},
	{ //34
		i_interval_string: "  |------========------|* ",
		x_interval_string: "  |------========------|* ",
		counter:           "34",
	},
	{ //35
		i_interval_string: "  |--------------======|* ",
		x_interval_string: "  |------========------|* ",
		counter:           "35",
	},
	{ //36
		i_interval_string: "  |=====---------------|  ",
		x_interval_string: " <|------=======-------|  ",
		counter:           "36",
	},
	{ //37
		i_interval_string: "  |======--------------|  ",
		x_interval_string: " <|------=======-------|  ",
		counter:           "37",
	},
	{ //38
		i_interval_string: "  |---=====------------|  ",
		x_interval_string: " <|------=======-------|  ",
		counter:           "38",
	},
	{ //39
		i_interval_string: "  |-------=====--------|  ",
		x_interval_string: " <|------========------|  ",
		counter:           "39",
	},
	{ //40
		i_interval_string: "  |------========------|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "40",
	},
	{ //41
		i_interval_string: "  |-----------=====----|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "41",
	},
	{ //42
		i_interval_string: "  |--------------======|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "42",
	},
	{ //43
		i_interval_string: "  |---------------=====|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "43",
	},
	{ //44
		i_interval_string: " <|-------=====--------|  ",
		x_interval_string: " <|------========------|  ",
		counter:           "44",
	},
	{ //45
		i_interval_string: " <|------========------|  ",
		x_interval_string: " <|------========------|  ",
		counter:           "45",
	},
	{ //46
		i_interval_string: "  |-------=====--------|> ",
		x_interval_string: "  |------========------|> ",
		counter:           "46",
	},
	{ //47
		i_interval_string: " <|------========------|> ",
		x_interval_string: " <|------========------|> ",
		counter:           "47",
	},
	{ //48
		i_interval_string: " <|--------====--------|> ",
		x_interval_string: " <|------========------|> ",
		counter:           "48",
	},
}

var testsParseInterval = []struct {
	s              string
	begin          int
	lowerIncluded  bool
	lowerUnbounded bool
	end            int
	upperIncluded  bool
	upperUnbounded bool
}{
	{
		s:              "  |=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " |=====|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	//-------------
	{
		s:              "  |--=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	//-------------
	{
		s:              "  |--=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====--|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====--|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====--|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====--|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	//-------------
	{
		s:              "  |=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====--|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====--|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====--|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====--|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
}

var testsIntervalSubtract = []struct {
	//i
	test testGeneral
	//g,h
	// i_interval_string.Subtract(x_interval_string)
	i_Subtract_x_before, i_Subtract_x_after string
	//j,k
	// x_interval_string.Subtract(i_interval_string)
	x_Subtract_i_before, x_Subtract_i_after string
}{
	{
		test:                testsGeneralSets[0],
		i_Subtract_x_before: "  |=====|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "|------=======-------|",
	},
	{
		test:                testsGeneralSets[1],
		i_Subtract_x_before: "  |======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|",
	},
	{
		test:                testsGeneralSets[2],
		i_Subtract_x_before: "  |---===|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|--------=====-----|",
	},
	{
		test:                testsGeneralSets[3],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "  |------=|*",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[4],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[5],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------==|",
		x_Subtract_i_before: "  |------=====|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[6],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "  |------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[7],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "  |---------------=====|",
		x_Subtract_i_before: "  |------========------|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[8],
		i_Subtract_x_before: " *|======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|",
	},
	{
		test:                testsGeneralSets[9],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "  |------&-|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[10],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "  |------========------|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[11],
		i_Subtract_x_before: "  |======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "  |------=======---|",
	},
	{
		test:                testsGeneralSets[12],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "  |--------------&|",
	},
	{
		test:                testsGeneralSets[13],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------======|*",
		x_Subtract_i_before: "  |------========------|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[14],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[15],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[16],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  " *|--------------==----|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[17],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[18],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[19],
		i_Subtract_x_before: " |======|*",
		i_Subtract_x_after:  "*|-------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[20],
		i_Subtract_x_before: " |======|*",
		i_Subtract_x_after:  "*|-------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[21],
		i_Subtract_x_before: " |---===|*",
		i_Subtract_x_after:  "*|-------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[22],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------&-------|>",
		x_Subtract_i_before: " |------=|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[23],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[24],
		i_Subtract_x_before: "  |======|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|",
	},
	{
		test:                testsGeneralSets[25],
		i_Subtract_x_before: "  |------&|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[26],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------======|",
		x_Subtract_i_before: "*|------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[27],
		i_Subtract_x_before: "|======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|*",
	},
	{
		test:                testsGeneralSets[28],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "|--------------&|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[29],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "|--------------======|",
		x_Subtract_i_before: " |------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[30],
		i_Subtract_x_before: "*|======|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "*|------=======|",
	},
	{
		test:                testsGeneralSets[31],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[32],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------======|",
		x_Subtract_i_before: "*|------========|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[33],
		i_Subtract_x_before: "  |======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "  |------=======|*",
	},
	{
		test:                testsGeneralSets[34],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[35],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "|--------------======|*",
		x_Subtract_i_before: "|------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[36],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "<|&|*",
		x_Subtract_i_after:  " *|-----========|",
	},
	{
		test:                testsGeneralSets[37],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "<|&|*",
		x_Subtract_i_after:  "*|------=======|",
	},
	{
		test:                testsGeneralSets[38],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "<|---&|*",
		x_Subtract_i_after:  "*|--------=====|",
	},
	{
		test:                testsGeneralSets[39],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: " <|------=|*",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[40],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "*|--------------&-------|>",
	},
	{
		test:                testsGeneralSets[41],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: " |------=====|*",
		x_Subtract_i_after:  "*|----------------&|>",
	},
	{
		test:                testsGeneralSets[42],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: " |------========|*",
		x_Subtract_i_after:  "*|--------------------&|>",
	},
	{
		test:                testsGeneralSets[43],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "|------=========------|*",
		x_Subtract_i_after:  "*|--------------------&|>",
	},
	{
		test:                testsGeneralSets[44],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[45],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[46],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "  |------=|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[47],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[48],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
}

var testsIntervalAdjoin = []struct {
	test testGeneral
	// i_interval_string.Adjoin(x_interval_string)
	i_Adjoin_x string
}{
	{
		test:       testsGeneralSets[0],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[1],
		i_Adjoin_x: "|=============|",
	},
	{
		test:       testsGeneralSets[2],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[3],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[4],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[5],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[6],
		i_Adjoin_x: "  |------==============|",
	},
	{
		test:       testsGeneralSets[7],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[8],
		i_Adjoin_x: "*|=============|",
	},
	{
		test:       testsGeneralSets[9],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[10],
		i_Adjoin_x: "|------==============|",
	},
	{
		test:       testsGeneralSets[11],
		i_Adjoin_x: "|=============|",
	},
	{
		test:       testsGeneralSets[12],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[13],
		i_Adjoin_x: "|------==============|*",
	},
	{
		test:       testsGeneralSets[14],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[15],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[16],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[17],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[18],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[19],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[20],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[21],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[22],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[23],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[24],
		i_Adjoin_x: "|=============|",
	},
	{
		test:       testsGeneralSets[25],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[26],
		i_Adjoin_x: "*|------==============|",
	},
	{
		test:       testsGeneralSets[27],
		i_Adjoin_x: "|=============|*",
	},
	{
		test:       testsGeneralSets[28],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[29],
		i_Adjoin_x: "|------==============|",
	},
	{
		test:       testsGeneralSets[30],
		i_Adjoin_x: "*|=============|",
	},
	{
		test:       testsGeneralSets[31],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[32],
		i_Adjoin_x: "*|------==============|",
	},
	{
		test:       testsGeneralSets[33],
		i_Adjoin_x: "|=============|*",
	},
	{
		test:       testsGeneralSets[34],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[35],
		i_Adjoin_x: "|------==============|*",
	},
	{
		test:       testsGeneralSets[36],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[37],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[38],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[39],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[40],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[41],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[42],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[43],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[44],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[45],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[46],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[47],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[48],
		i_Adjoin_x: "",
	},
}

var testsIntervalContains = []struct {
	test      testGeneral
	i_Cover_x bool
	x_Cover_i bool
}{
	{
		test:      testsGeneralSets[0],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[1],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[2],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[3],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[4],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[5],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[6],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[7],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[8],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[9],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[10],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[11],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[12],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[13],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[14],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[15],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[16],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[17],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[18],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[19],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[20],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[21],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[22],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[23],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[24],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[25],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[26],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[27],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[28],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[29],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[30],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[31],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[32],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[33],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[34],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[35],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[36],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[37],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[38],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[39],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[40],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[41],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[42],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[43],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[44],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[45],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[46],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[47],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[48],
		i_Cover_x: true,
		x_Cover_i: true,
	},
}

var testsIntervalLeEndOf = []struct {
	test      testGeneral
	i_LeEnd_x bool
	x_LeEnd_i bool
}{
	{
		test:      testsGeneralSets[0],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[1],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[2],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[3],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[4],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[5],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[6],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[7],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[8],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[9],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[10],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[11],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[12],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[13],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[14],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[15],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[16],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[17],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[18],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[19],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[20],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[21],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[22],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[23],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[24],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[25],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[26],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[27],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[28],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[29],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[30],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[31],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[32],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[33],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[34],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[35],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[36],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[37],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[38],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[39],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[40],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[41],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[42],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[43],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[44],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[45],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[46],
		i_LeEnd_x: false,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[47],
		i_LeEnd_x: false,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[48],
		i_LeEnd_x: false,
		x_LeEnd_i: false,
	},
}

var testsIntervalLtBeginOf = []struct {
	test       testGeneral
	i_Before_x bool
	x_Before_i bool
}{
	{
		test:       testsGeneralSets[0],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[1],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[2],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[3],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[4],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[5],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[6],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[7],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[8],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[9],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[10],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[11],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[12],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[13],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[14],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[15],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[16],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[17],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[18],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[19],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[20],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[21],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[22],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[23],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[24],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[25],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[26],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[27],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[28],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[29],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[30],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[31],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[32],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[33],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[34],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[35],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[36],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[37],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[38],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[39],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[40],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[41],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[42],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[43],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[44],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[45],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[46],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[47],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[48],
		i_Before_x: false,
		x_Before_i: false,
	},
}

var testsIntervalIntersect = []struct {
	test          testGeneral
	i_intersect_x string
}{
	{
		test:          testsGeneralSets[0],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[1],
		i_intersect_x: "|------&------|",
	},
	{
		test:          testsGeneralSets[2],
		i_intersect_x: "|------==------|",
	},
	{
		test:          testsGeneralSets[3],
		i_intersect_x: "|-------=====------|",
	},
	{
		test:          testsGeneralSets[4],
		i_intersect_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[5],
		i_intersect_x: "|-----------===------|",
	},
	{
		test:          testsGeneralSets[6],
		i_intersect_x: "|--------------&------|",
	},
	{
		test:          testsGeneralSets[7],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[8],
		i_intersect_x: "|------&------|",
	},
	{
		test:          testsGeneralSets[9],
		i_intersect_x: "*|------========------|",
	},
	{
		test:          testsGeneralSets[10],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[11],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[12],
		i_intersect_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[13],
		i_intersect_x: "|--------------&|",
	},
	{
		test:          testsGeneralSets[14],
		i_intersect_x: "|------======|",
	},
	{
		test:          testsGeneralSets[15],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[16],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[17],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[18],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[19],
		i_intersect_x: "|------=======|",
	},
	{
		test:          testsGeneralSets[20],
		i_intersect_x: "|------=======|",
	},
	{
		test:          testsGeneralSets[21],
		i_intersect_x: "|------=======|",
	},
	{
		test:          testsGeneralSets[22],
		i_intersect_x: "|-------=======|",
	},
	{
		test:          testsGeneralSets[23],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[24],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[25],
		i_intersect_x: "*|------========|",
	},
	{
		test:          testsGeneralSets[26],
		i_intersect_x: "|--------------&|",
	},
	{
		test:          testsGeneralSets[27],
		i_intersect_x: "|------&|",
	},
	{
		test:          testsGeneralSets[28],
		i_intersect_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[29],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[30],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[31],
		i_intersect_x: "*|------========|",
	},
	{
		test:          testsGeneralSets[32],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[33],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[34],
		i_intersect_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[35],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[36],
		i_intersect_x: "|=====----|",
	},
	{
		test:          testsGeneralSets[37],
		i_intersect_x: "|======----|",
	},
	{
		test:          testsGeneralSets[38],
		i_intersect_x: "|---=====------------|",
	},
	{
		test:          testsGeneralSets[39],
		i_intersect_x: "|-------=====--------|",
	},
	{
		test:          testsGeneralSets[40],
		i_intersect_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[41],
		i_intersect_x: "|-----------=====----|",
	},
	{
		test:          testsGeneralSets[42],
		i_intersect_x: "|--------------======|",
	},
	{
		test:          testsGeneralSets[43],
		i_intersect_x: "|---------------=====|",
	},
	{
		test:          testsGeneralSets[44],
		i_intersect_x: "<|------======------|",
	},
	{
		test:          testsGeneralSets[45],
		i_intersect_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[46],
		i_intersect_x: "|-------=====--------|>",
	},
	{
		test:          testsGeneralSets[47],
		i_intersect_x: "<|------========------|>",
	},
	{
		test:          testsGeneralSets[48],
		i_intersect_x: "<|------========------|>",
	},
}

var testsIntervalEncompass = []struct {
	// i_interval_string.Encompass(x_interval_string)
	i_Encompass_x string
	test          testGeneral
}{
	{
		test:          testsGeneralSets[0],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[1],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[2],
		i_Encompass_x: "|---==========|",
	},
	{
		test:          testsGeneralSets[3],
		i_Encompass_x: "|------========|",
	},
	{
		test:          testsGeneralSets[4],
		i_Encompass_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[5],
		i_Encompass_x: "|------==========------|",
	},
	{
		test:          testsGeneralSets[6],
		i_Encompass_x: "|------==============------|",
	},
	{
		test:          testsGeneralSets[7],
		i_Encompass_x: "|------==============------|",
	},
	{
		test:          testsGeneralSets[8],
		i_Encompass_x: "*|=============|",
	},
	{
		test:          testsGeneralSets[9],
		i_Encompass_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[10],
		i_Encompass_x: "|------==============------|",
	},
	{
		test:          testsGeneralSets[11],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[12],
		i_Encompass_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[13],
		i_Encompass_x: "|------==============------|*",
	},
	{
		test:          testsGeneralSets[14],
		i_Encompass_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[15],
		i_Encompass_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[16],
		i_Encompass_x: "<|------==========------|",
	},
	{
		test:          testsGeneralSets[17],
		i_Encompass_x: "<|------==============------|",
	},
	{
		test:          testsGeneralSets[18],
		i_Encompass_x: "<|------==============------|",
	},
	{
		test:          testsGeneralSets[19],
		i_Encompass_x: "|=============------|>",
	},
	{
		test:          testsGeneralSets[20],
		i_Encompass_x: "|=============|>", //0 13
	},
	{
		test:          testsGeneralSets[21],
		i_Encompass_x: "|---==========|>",
	},
	{
		test:          testsGeneralSets[22],
		i_Encompass_x: "|------========|>", //6 14
	},
	{
		test:          testsGeneralSets[23],
		i_Encompass_x: "|------========|>",
	},
	{
		test:          testsGeneralSets[24],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[25],
		i_Encompass_x: "|------========|",
	},
	{
		test:          testsGeneralSets[26],
		i_Encompass_x: "*|------==============--|", //6 20
	},
	{
		test:          testsGeneralSets[27],
		i_Encompass_x: "|=============|*",
	},
	{
		test:          testsGeneralSets[28],
		i_Encompass_x: "|------========|",
	},
	{
		test:          testsGeneralSets[29],
		i_Encompass_x: "|------==============--|",
	},
	{
		test:          testsGeneralSets[30],
		i_Encompass_x: "*|=============|",
	},
	{
		test:          testsGeneralSets[31],
		i_Encompass_x: "*|------========|",
	},
	{
		test:          testsGeneralSets[32],
		i_Encompass_x: "*|------==============--|",
	},
	{
		test:          testsGeneralSets[33],
		i_Encompass_x: "|=============|*",
	},
	{
		test:          testsGeneralSets[34],
		i_Encompass_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[35],
		i_Encompass_x: "|------==============--|*",
	},
	{
		test:          testsGeneralSets[36],
		i_Encompass_x: "<|=============|",
	},
	{
		test:          testsGeneralSets[37],
		i_Encompass_x: "<|=============|",
	},
	{
		test:          testsGeneralSets[38],
		i_Encompass_x: "<|---==========|",
	},
	{
		test:          testsGeneralSets[39],
		i_Encompass_x: "<|------========|",
	},
	{
		test:          testsGeneralSets[40],
		i_Encompass_x: "|------========|>",
	},
	{
		test:          testsGeneralSets[41],
		i_Encompass_x: "|------==========|>",
	},
	{
		test:          testsGeneralSets[42],
		i_Encompass_x: "|------==============|>",
	},
	{
		test:          testsGeneralSets[43],
		i_Encompass_x: "|------==============|>",
	},
	{
		test:          testsGeneralSets[44],
		i_Encompass_x: "<|------========|",
	},
	{
		test:          testsGeneralSets[45],
		i_Encompass_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[46],
		i_Encompass_x: "|------========|>",
	},
	{
		test:          testsGeneralSets[47],
		i_Encompass_x: "<|------========------|>",
	},
	{
		test:          testsGeneralSets[48],
		i_Encompass_x: "<|------========|>",
	},
}

var testsHAS = []struct {
	s       string
	value   int
	result  bool
	counter string
}{
	{
		s:       "  |===|",
		value:   3,
		result:  true,
		counter: "0",
	},
	{
		s:       "  |===|",
		value:   0,
		result:  true,
		counter: "1",
	},
	{
		s:       "  |===|",
		value:   -1,
		result:  false,
		counter: "2",
	},
	{
		s:       "  |===|",
		value:   6,
		result:  false,
		counter: "3",
	},
	//----------------------------
	{
		s:       "*|===|",
		value:   3,
		result:  true,
		counter: "4",
	},
	{
		s:       "*|===|",
		value:   0,
		result:  false,
		counter: "5",
	},
	{
		s:       "*|===|",
		value:   -1,
		result:  false,
		counter: "6",
	},
	{
		s:       "*|===|",
		value:   6,
		result:  false,
		counter: "7",
	},
	//------------------------
	{
		s:       "|=====|*",
		value:   5,
		result:  false,
		counter: "8",
	},
	{
		s:       "|=====|*",
		value:   0,
		result:  true,
		counter: "9",
	},
	{
		s:       "|=====|*",
		value:   -1,
		result:  false,
		counter: "10",
	},
	{
		s:       "|=====|*",
		value:   6,
		result:  false,
		counter: "11",
	},
	//------------------------
	{
		s:       "*|=====|*",
		value:   5,
		result:  false,
		counter: "12",
	},
	{
		s:       "*|=====|*",
		value:   0,
		result:  false,
		counter: "13",
	},
	{
		s:       "*|=====|*",
		value:   -1,
		result:  false,
		counter: "14",
	},
	{
		s:       "*|=====|*",
		value:   6,
		result:  false,
		counter: "15",
	},
	//----------------------------
	{
		s:       " <|=====|",
		value:   5,
		result:  true,
		counter: "16",
	},
	{
		s:       " <|=====|",
		value:   0,
		result:  true,
		counter: "17",
	},
	{
		s:       " <|=====|",
		value:   -1,
		result:  true,
		counter: "18",
	},
	{
		s:       " <|=====|",
		value:   6,
		result:  false,
		counter: "19",
	},
	//----------------------------
	{
		s:       " |=====|>",
		value:   5,
		result:  true,
		counter: "20",
	},
	{
		s:       " |=====|>",
		value:   0,
		result:  true,
		counter: "21",
	},
	{
		s:       " |=====|>",
		value:   -1,
		result:  false,
		counter: "22",
	},
	{
		s:       " |=====|>",
		value:   6,
		result:  true,
		counter: "23",
	},
	//----------------------------
	{
		s:       " <|=====|>",
		value:   5,
		result:  true,
		counter: "24",
	},
	{
		s:       " <|=====|>",
		value:   0,
		result:  true,
		counter: "25",
	},
	{
		s:       " <|=====|>",
		value:   -1,
		result:  true,
		counter: "26",
	},
	{
		s:       " <|=====|>",
		value:   6,
		result:  true,
		counter: "27",
	},
	//----------------------------
	{
		s:       "  |=====|*>",
		value:   5,
		result:  true,
		counter: "28",
	},
	{
		s:       "  |=====|*>",
		value:   0,
		result:  true,
		counter: "29",
	},
	{
		s:       "  |=====|*>",
		value:   -1,
		result:  false,
		counter: "30",
	},
	{
		s:       "  |=====|*>",
		value:   6,
		result:  true,
		counter: "31",
	},
	//----------------------------
	{
		s:       "<*|=====|",
		value:   5,
		result:  true,
		counter: "32",
	},
	{
		s:       "<*|=====|",
		value:   0,
		result:  true,
		counter: "33",
	},
	{
		s:       "<*|=====|",
		value:   -1,
		result:  true,
		counter: "34",
	},
	{
		s:       "<*|=====|",
		value:   6,
		result:  false,
		counter: "35",
	},
	//----------------------------
	{
		s:       "<*|=====|*>",
		value:   5,
		result:  true,
		counter: "36",
	},
	{
		s:       "<*|=====|*>",
		value:   0,
		result:  true,
		counter: "37",
	},
	{
		s:       "<*|=====|*>",
		value:   -1,
		result:  true,
		counter: "38",
	},
	{
		s:       "<*|=====|*>",
		value:   6,
		result:  true,
		counter: "39",
	},
	//----------------------------
}
