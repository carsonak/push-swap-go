#!/usr/bin/env python3
import argparse
import random

parser = argparse.ArgumentParser(
    description="Generate random numbers and print them on one line."
)

parser.add_argument("count", type=int, help="number of random numbers to generate")
parser.add_argument("min", type=int, help="minimum value (inclusive)")
parser.add_argument("max", type=int, help="maximum value (inclusive)")

parser.add_argument(
    "--allow-duplicates",
    action="store_true",
    help="allow duplicate numbers (default: duplicates not allowed)",
)

parser.add_argument(
    "--seed",
    type=float,
    help="seed for reproducible output",
)

args = parser.parse_args()

if args.min > args.max:
    parser.error("min must be less than or equal to max")

if args.seed is not None:
    random.seed(args.seed)

range_size = args.max - args.min + 1

if args.allow_duplicates:
    numbers = [random.uniform(args.min, args.max) for _ in range(args.count)]
else:
    if args.count > range_size:
        parser.error("count is larger than the number of unique values in the range")

    numbers = random.sample(range(args.min, args.max + 1), args.count)

print(" ".join(str(n) for n in numbers))
