package main

var htmlTemplate string = `
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
    <title>Squirrel Library: {{}}</title>

	<link rel="stylesheet" media="screen" href="style.css">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="author" content="James Kilfiger">
	<meta name="description" content="A mysql library">
	<meta name="robots" content="all">


	<!--[if lt IE 9]>
	<script src="script/html5shiv.js"></script>
	<![endif]-->
</head>


<body id="css-zen-garden">
<div class="page-wrapper">

	<section class="intro" id="zen-intro">
		<header role="banner">
            <h1>Squirrel Library</h1>
            <h2>{{}}</h2>
		</header>

		<div class="summary" id="zen-summary" role="article">
            {{}}
		</div>

		<div class="preamble" id="zen-preamble" role="article">
            {{}}
		</div>
	</section>

	<div class="main supporting" id="zen-supporting" role="main">
		<div class="explanation" id="zen-explanation" role="article">
		</div>

		<div class="participation" id="zen-participation" role="article">
		</div>

		<div class="benefits" id="zen-benefits" role="article">
		</div>

		<div class="requirements" id="zen-requirements" role="article">
		</div>

		<footer>
		</footer>

	</div>


	<aside class="sidebar" role="complementary">
		<div class="wrapper">

			<div class="design-selection" id="design-selection">
				<h3 class="select">Select a Design:</h3>
				<nav role="navigation">
				</nav>
			</div>

			<div class="design-archives" id="design-archives">
				<h3 class="archives">Archives:</h3>
				<nav role="navigation">
				</nav>
			</div>

			<div class="zen-resources" id="zen-resources">
				<h3 class="resources">Resources:</h3>
			</div>
		</div>
	</aside>


</div>

</body>
</html>
`
